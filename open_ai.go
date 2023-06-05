package main
import (
"fmt"
"net/http"
"github.com/gorilla/websocket"
"crypto/tls"
"io/ioutil"
"bytes"
"github.com/tidwall/gjson"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 409600,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func get_answer_from_openai(question string) (answer string) {
    var req *http.Request
    var err error
    question_json := "{\"model\":\"gpt-3.5-turbo\", \"messages\":[{\"role\":\"user\", \"content\":\"" + question + "\"}]}"
    fmt.Println(question_json)
    if req, err = http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer([]byte(question_json))); err != nil {
        fmt.Println("new request error")
        return
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer")

    tr := &http.Transport{TLSClientConfig:&tls.Config{InsecureSkipVerify:true}}
    client := &http.Client{Transport:tr}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

    answer_, _ := ioutil.ReadAll(resp.Body)
    value := string(answer_)
    array := gjson.Get(value, "choices")
    for _, v := range array.Array() {
        array_ := gjson.Get(v.String(), "message.content")
        answer = string(array_.String())
	}
    return
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
        answer := get_answer_from_openai(question)
        if err := conn.WriteMessage(messageType, []byte(answer)); err != nil {
	        fmt.Println(err)
            return
        }
    }
}

func main() {
    // res := get_answer_from_openai("who are you?")
	// fmt.Println(res)
    // return
    http.HandleFunc("/openai", openai)
    http.ListenAndServe(":5757", nil)
}
