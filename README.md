# Simple Golang Chat Application

A lightweight, real-time chat application built with Go that demonstrates TCP networking, concurrent programming, and client-server architecture.

The project includes a web-based group chat and the original command-line TCP chat. The web interface can open multiple room windows on one page, and everyone in a room can chat in real time.

## 🚀 Features

- **Real-time messaging**: Send and receive messages instantly through a central server
- **Multi-client support**: Connect multiple clients simultaneously
- **Simple command-line interface**: Easy-to-use terminal-based chat
- **Web chat interface**: Open multiple chat room windows in one page
- **Group rooms**: Everyone in the same room receives messages in real time
- **Concurrent processing**: Server handles multiple connections using goroutines
- **JSON message format**: Structured message passing with timestamps

## 🏗️ Architecture

The application follows a client-server architecture:

- **Server**: Central hub that manages client connections and routes messages
- **Client**: Terminal-based chat client that connects to the server
- **Message**: Structured data format for communication

```
Client A ←→ Server ←→ Client B
Client C ←→     ←→ Client D
```

## 📋 Prerequisites

- Go 1.20 or higher
- Terminal/Command prompt

## 🛠️ Installation & Usage

### 1. Start the Server

From the project root, start the web and TCP servers:

```bash
go run ./server
```

Open the web chat at [http://localhost:8083](http://localhost:8083). Browse or search the room list, select a room, then click **Open Chat Window** and enter a nickname for that window. Click **Create Room** to add a room, then select it from the list. People in the same room can chat together.

The TCP backend listens on port `8082`; the web frontend listens on port `8083`.


### Start a Command-line Client (optional)

In another terminal, run:

```bash
go run ./client
```

Each client will be assigned a unique ID by the server. You'll see output like:
```
Assign id 0 to 127.0.0.1:xxxxx
```

### Send Messages

To send a message to another client, use the following command format:

```
-t <recipient_id> -m <your_message>
```

**Example:**
```
-t 1 -m Hello, how are you?
-t 0 -m This is a test message
```

### Exit the Application

Type `exit` in any client to disconnect from the server.

## 📝 Message Format

Messages are structured as JSON with the following fields:

```json
{
  "From": "sender_id",
  "To": "recipient_id", 
  "Text": "message_content",
  "CreatedAt": "2024-01-01T12:00:00Z"
}
```

## 🔧 Technical Details

- **Protocol**: TCP
- **TCP backend port**: 8082
- **Web frontend port**: 8083
- **Message Format**: JSON
- **Concurrency**: Goroutines for handling multiple connections
- **Buffer Size**: 1024 bytes per message

## 🎯 Use Cases

- Learning Go networking concepts
- Understanding client-server architecture
- Demonstrating concurrent programming with goroutines
- Simple team communication tool
- Educational project for TCP socket programming

## 📁 Project Structure

```
simple-chat-application/
├── server/
│   ├── server.go          # TCP server and application entry point
│   ├── web.go             # Web chat HTTP API and real-time event stream
│   └── index.html         # Web chat interface
├── client/
│   └── client.go          # Client implementation
├── message/
│   └── message.go         # Message structure definition
└── README.md              # This file
```

## 🤝 Contributing

Feel free to submit issues and enhancement requests!

## 📄 License

This project is open source and available under the [MIT License](LICENSE).
