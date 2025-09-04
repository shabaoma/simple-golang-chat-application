# Simple Golang Chat Application

A lightweight, real-time chat application built with Go that demonstrates TCP networking, concurrent programming, and client-server architecture.

## 🚀 Features

- **Real-time messaging**: Send and receive messages instantly through a central server
- **Multi-client support**: Connect multiple clients simultaneously
- **Simple command-line interface**: Easy-to-use terminal-based chat
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

- Go 1.11 or higher
- Terminal/Command prompt

## 🛠️ Installation & Usage

### 1. Start the Server

First, start the central server that will handle all client connections:

```bash
cd server
go run server.go
```

The server will start listening on port `8080` and display:
```
Start server
Wait for client to connect
```

### 2. Start Multiple Clients

Open new terminal windows and start as many clients as needed:

```bash
cd client
go run client.go
```

Each client will be assigned a unique ID by the server. You'll see output like:
```
Assign id 0 to 127.0.0.1:xxxxx
```

### 3. Send Messages

To send a message to another client, use the following command format:

```
-t <recipient_id> -m <your_message>
```

**Example:**
```
-t 1 -m Hello, how are you?
-t 0 -m This is a test message
```

### 4. Exit the Application

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
- **Port**: 8080
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
│   └── server.go          # Central server implementation
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
