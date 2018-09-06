package main

import (
	"log"
	"net/http"
	"strings"
)

func wasmHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".wasm") {
		w.Header().Set("Content-Type", "application/wasm")
		log.Printf("wasm path is /%s\n", r.URL.Path[1:])
	} else {
		log.Printf("path is /%s\n", r.URL.Path[1:])
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", wasmHandler)
	log.Println("Listening on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
