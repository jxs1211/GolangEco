package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG\n")
	time.Sleep(2 * time.Second)
}

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	http.HandleFunc("/ping", http.HandlerFunc(ping))
	log.Fatal(http.ListenAndServe(":48080", nil))
}
