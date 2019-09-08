package main

import (
        "log"
        "net/http"
        "github.com/p0dxd/subway/lib"
)


func main() {
        subway.Init()
        http.Handle("/", http.FileServer(http.Dir("./static")))
        log.Print("Now serving: http://localhost:3001")
        http.ListenAndServe(":3001", nil)
}
