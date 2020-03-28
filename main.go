package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	Libs "./pkg/libs"
)

func main() {
	initMain()
	Libs.Init()
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Print("Now serving: https://localhost:3001")
	http.ListenAndServe(":3001", nil)

	err := http.ListenAndServeTLS(":3001", "fullchain.pem", "privkey.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initMain() {
	writePublicKey()
	writePrivateKey()
}

func writePublicKey() {
	publicKey := os.Getenv("PUBLIC_KEY")
	data, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Fatal("error:", err)
	}

	f, err := os.Create("fullchain.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.Write(data)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func writePrivateKey() {
	privateKey := os.Getenv("PRIVATE_KEY")

	data, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		log.Fatal("error:", err)
	}

	f, err := os.Create("privkey.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.Write(data)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
