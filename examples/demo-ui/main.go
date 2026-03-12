package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("DEMO_UI_PORT")
	if port == "" {
		port = "3000"
	}

	// Serve the static folder containing index.html
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	log.Printf("Starting Gantral Demo UI on :%s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
