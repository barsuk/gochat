package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

var storage []string

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/chatserver", chatserver)
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func chat(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "chat.html")
}

func home(w http.ResponseWriter, req *http.Request) {
	//io.WriteString(w, "<h2>hola el mundo</h2>")
	http.ServeFile(w, req, "index.html")
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("on upgrade: ", err)
		return
	}
	defer c.Close()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Print("не прочитал: ", err)
			break
		}
		log.Printf("Я принял: %s", message)
		err = c.WriteMessage(messageType, message)
		if err != nil {
			log.Println("не ответил: ", err)
			break
		}
	}
}

func chatserver(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("on upgrade: ", err)
		return
	}
	defer c.Close()

	// отдаем содержимое чата
	go func() {
		var storageLength = 0
		for {
			if len(storage) > storageLength {
				for _, msg := range storage[storageLength:] {
					err = c.WriteMessage(1, []byte(fmt.Sprintf("%s", msg)))
					if err != nil {
						log.Printf("на msg: %s не ответил: %s", msg, err)
					}
				}
			}

			storageLength = len(storage)

			time.Sleep(1)
		}
	}()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Print("не прочитал: ", err)
			break
		}
		log.Printf("Я принял: %d ### %s", messageType, message)
		storage = append(storage, string(message))
	}
}
