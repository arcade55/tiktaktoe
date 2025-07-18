package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/starfederation/datastar/sdk/go/datastar"
)

type Signals struct {
	Cell  string `json:"cell"`
	Shape string `json:"shape"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, homepage().Render())
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/cell", func(w http.ResponseWriter, r *http.Request) {
		signals := Signals{}
		if err := datastar.ReadSignals(r, &signals); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("cell = %s shape = %s", signals.Cell, signals.Shape)

		broker.broadcast(button(signals).Render())
	})

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		broker.broadcast(board().Render())
	})

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {

		signals := Signals{}
		sse := datastar.NewSSE(w, r)
		broker.mu.Lock()
		broker.clients[sse] = struct{}{}

		if len(broker.clients) == 1 {
			signals.Shape = "X"
		} else {
			signals.Shape = "O"
		}

		sse.MarshalAndPatchSignals(signals)
		log.Printf("Client connected, total: %d", len(broker.clients))
		broker.mu.Unlock()

		<-r.Context().Done()

		broker.mu.Lock()
		delete(broker.clients, sse)
		log.Printf("Client disconnected, total: %d", len(broker.clients))
		broker.mu.Unlock()
	})
	// Start the server on port 8080.
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
