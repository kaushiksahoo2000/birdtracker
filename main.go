package main

import (
	"crypto/rand"
	"fmt"
	"log"
	random "math/rand"
	"time"
)

// The PRODUCER INTERVAL is how often
// a BirdTrack will be produced
const PRODUCER_INTERVAL = 1 * time.Second

// BirdTrack represents a fake Bird vehicle,
// and contains velocity information
type BirdTrack struct {
	ID        string
	Latitude  float64
	Longitude float64
	Speed     int
}

// GenUUID is a helper function to generate a
// UUID string for each BirdTrack that is created
func GenUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

// NewBirdTrack generates a new BirdTrack struct
// with random values for velocity information
func NewBirdTrack() *BirdTrack {
	return &BirdTrack{
		ID:        GenUUID(),
		Latitude:  random.Float64() * 90,
		Longitude: random.Float64()*360 - 180,
		Speed:     random.Intn(19), // 19 mph seems to be the max speed for Birds
	}
}

// NewBirdTrackProducer creates and returns a channel
// that it sends a new BirdTrack on every PRODUCER_INTERVAL (1 second currently)
func NewBirdTrackProducer() <-chan *BirdTrack {
	producer := make(chan *BirdTrack)
	ticker := time.NewTicker(PRODUCER_INTERVAL)
	go func() {
		defer func() {
			close(producer)
			ticker.Stop()
		}()
		for {
			select {
			case <-ticker.C:
				producer <- NewBirdTrack()
			}
		}
	}()
	return producer
}

func main() {
	producer := NewBirdTrackProducer()
	// This for block is a consumer
	// for the producer we just created
	for {
		birdTrack := <-producer
		if birdTrack.Speed > 10 {
			fmt.Println("Bird is going faster than 10 mph", "ID: ", birdTrack.ID, "Latitude: ", birdTrack.Latitude, "Longitude: ", birdTrack.Longitude, "Speed: ", birdTrack.Speed)
		}
	}
}
