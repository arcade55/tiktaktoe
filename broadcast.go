package main

import (
	"log"
	"sync"

	"github.com/starfederation/datastar/sdk/go/datastar"
)

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
