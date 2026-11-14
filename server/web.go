package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

//go:embed index.html
var webFiles embed.FS

type webMessage struct {
	From      string    `json:"from"`
	Room      string    `json:"room"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type webHub struct {
	mu      sync.Mutex
	clients map[chan webMessage]struct{}
	history []webMessage
	rooms   map[string]struct{}
}

var chatHub = &webHub{
	clients: make(map[chan webMessage]struct{}),
	rooms:   map[string]struct{}{"大厅": {}},
}

func (h *webHub) addRoom(room string) {
	h.mu.Lock()
	h.rooms[room] = struct{}{}
	h.mu.Unlock()
}

func (h *webHub) addClient(ch chan webMessage) {
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
}

func (h *webHub) removeClient(ch chan webMessage) {
	h.mu.Lock()
	delete(h.clients, ch)
	close(ch)
	h.mu.Unlock()
}

func (h *webHub) publish(m webMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.history = append(h.history, m)
	h.rooms[m.Room] = struct{}{}
	if len(h.history) > 200 {
		h.history = h.history[len(h.history)-200:]
	}
	for ch := range h.clients {
		select {
		case ch <- m:
		default:
			// Keep the newest messages flowing if a browser is temporarily slow.
			select {
			case <-ch:
			default:
			}
			ch <- m
		}
	}
}

func startWebServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data, err := webFiles.ReadFile("index.html")
		if err != nil {
			http.Error(w, "网页文件暂不可用", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
	})
	mux.HandleFunc("/api/messages", handleMessages)
	mux.HandleFunc("/api/rooms", handleRooms)
	mux.HandleFunc("/api/events", handleEvents)
	log.Println("Web chat available at http://localhost:8083")
	if err := http.ListenAndServe(":8083", mux); err != nil {
		log.Println("Web server error:", err)
	}
}

func handleRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch r.Method {
	case http.MethodGet:
		chatHub.mu.Lock()
		rooms := make([]string, 0, len(chatHub.rooms))
		for room := range chatHub.rooms {
			rooms = append(rooms, room)
		}
		chatHub.mu.Unlock()
		sort.Strings(rooms)
		_ = json.NewEncoder(w).Encode(rooms)
	case http.MethodPost:
		var input struct {
			Room string `json:"room"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024)).Decode(&input); err != nil {
			http.Error(w, "无效的房间名称", http.StatusBadRequest)
			return
		}
		room := normalizeRoom(input.Room)
		if strings.TrimSpace(input.Room) == "" {
			http.Error(w, "房间名称不能为空", http.StatusBadRequest)
			return
		}
		chatHub.addRoom(room)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(room)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "仅支持 GET 和 POST", http.StatusMethodNotAllowed)
	}
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch r.Method {
	case http.MethodGet:
		room := normalizeRoom(r.URL.Query().Get("room"))
		chatHub.mu.Lock()
		messages := make([]webMessage, 0)
		for _, m := range chatHub.history {
			if m.Room == room {
				messages = append(messages, m)
			}
		}
		chatHub.mu.Unlock()
		_ = json.NewEncoder(w).Encode(messages)
	case http.MethodPost:
		var input struct {
			From string `json:"from"`
			Room string `json:"room"`
			Text string `json:"text"`
		}
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16*1024))
		if err := decoder.Decode(&input); err != nil && err != io.EOF {
			http.Error(w, "无效的消息数据", http.StatusBadRequest)
			return
		}
		input.From = strings.TrimSpace(input.From)
		input.Text = strings.TrimSpace(input.Text)
		if input.From == "" || input.Text == "" {
			http.Error(w, "昵称和消息不能为空", http.StatusBadRequest)
			return
		}
		if len([]rune(input.From)) > 32 || len([]rune(input.Text)) > 2000 {
			http.Error(w, "昵称最多 32 个字，消息最多 2000 个字", http.StatusBadRequest)
			return
		}
		m := webMessage{From: input.From, Room: normalizeRoom(input.Room), Text: input.Text, CreatedAt: time.Now()}
		chatHub.publish(m)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(m)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "仅支持 GET 和 POST", http.StatusMethodNotAllowed)
	}
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "当前服务器不支持实时推送", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	ch := make(chan webMessage, 32)
	chatHub.addClient(ch)
	defer chatHub.removeClient(ch)
	_, _ = fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()
	for {
		select {
		case <-r.Context().Done():
			return
		case m := <-ch:
			data, err := json.Marshal(m)
			if err != nil {
				continue
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				return
			}
			flusher.Flush()
		case <-time.After(20 * time.Second):
			_, _ = fmt.Fprint(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

func normalizeRoom(room string) string {
	room = strings.TrimSpace(room)
	if room == "" {
		return "大厅"
	}
	if len([]rune(room)) > 40 {
		return string([]rune(room)[:40])
	}
	return room
}
