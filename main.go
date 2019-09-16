package main

import (
        "log"
        "net/http"
        Libs "./pkg/libs"
)


func main() {
        Libs.Init()
        http.Handle("/", http.FileServer(http.Dir("./static")))
        log.Print("Now serving: http://localhost:3001")
        http.ListenAndServe(":3001", nil)
}
