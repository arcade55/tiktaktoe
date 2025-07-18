package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/starfederation/datastar/sdk/go/datastar"
)

type Signals struct {
	Cell  string `json:"cell"`
	Shape string `json:"shape"`
}

var broker = NewBroker()

type Broker struct {
	mu      sync.Mutex
	clients map[*datastar.ServerSentEventGenerator]struct{}
}

func NewBroker() *Broker {
	b := &Broker{
		clients: make(map[*datastar.ServerSentEventGenerator]struct{}),
	}
	return b
}

func (b *Broker) broadcast(elementToPatch string) {
	b.mu.Lock()
	toRemove := []*datastar.ServerSentEventGenerator{}
	for sse := range b.clients {
		if err := sse.PatchElements(elementToPatch); err != nil {
			log.Printf("Error sending to client: %v", err)
			toRemove = append(toRemove, sse)
		}
	}
	for _, sse := range toRemove {
		delete(b.clients, sse)
	}
	log.Printf("Broadcasted to %d clients", len(b.clients))
	b.mu.Unlock()
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
		broker.broadcast(fmt.Sprintf(`<button id="%s" class="cell">%s</button>`, signals.Cell, signals.Shape))
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
