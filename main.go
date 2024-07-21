package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var client = make(map[string]*Client)
var wait = make(chan bool)

type Client struct {
	id   string
	conn *websocket.Conn
	send chan []byte
}

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "chat.html")
	})
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	<-wait
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		conn.Close()
		return
	}
	c := &Client{id: id, conn: conn, send: make(chan []byte, 256)}
	client[id] = c
	go StartChat(id)
}
func StartChat(id string) {
	for {
		// read in a message
		messageType, p, err := client[id].conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))
		for _, c := range client {
			if c != client[id] {
				if err := c.conn.WriteMessage(messageType, p); err != nil {
					log.Println(err)
					return
				}
			}

		}
		/*if err := client[id].conn.WriteMessage(messageType, []byte("hi dost"+id)); err != nil {
			log.Println(err)
			return
		}*/
	}
}
