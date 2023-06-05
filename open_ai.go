package main
import (
"fmt"
"net/http"
"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 409600,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
func openai(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
	    fmt.Println(err)
        return
    }
    defer conn.Close()
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
	        fmt.Println(err)
            return
        }
        question := string(p[:])
	    fmt.Println(question)
        answer := "answer from server"
        if err := conn.WriteMessage(messageType, []byte(answer)); err != nil {
	        fmt.Println(err)
            return
        }
    }
}

func main() {
	fmt.Println("Hello, 世界")
    http.HandleFunc("/openai", openai)
    http.ListenAndServe(":5757", nil)
}
