package main

import (
	"fmt"
	"log"
	"net/http"
)

func wasmHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	fmt.Println("path is", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/examples/wasm/", wasmHandler)
	log.Fatal(http.ListenAndServe(":8081", mux))
}
