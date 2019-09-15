package main

import (
        "log"
        "net/http"
        L "./pkg/libs"
)


func main() {
        log.Print("demo")
        L.Init()
        http.Handle("/", http.FileServer(http.Dir("./static")))
        log.Print("Now serving: http://localhost:3001")
        http.ListenAndServe(":3001", nil)
}
