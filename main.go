package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/willeslau/kafka-producer/eventqueue"
)

func main() {
	broker := os.Args[1]
	topic := os.Args[2]

	config := eventqueue.KafkaConfig{broker}
	producer, err := eventqueue.NewProducer(&config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	http.HandleFunc("/produce", getHandler(producer, topic))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// Payload the payload
type Payload struct {
	Message string
}

func getHandler(p eventqueue.Producer, topic string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		payload := Payload{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, "Error!")
			return
		}

		p.Produce(topic, payload.Message)
		fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	}
	return fn
}
