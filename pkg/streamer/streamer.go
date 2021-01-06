// Package streamer replaces the original pubsub mechanism.
package streamer

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/sse"
)

var streamers = make(map[uint]*sse.Streamer)

// Publish send content to the stream
func Publish(id uint, content []byte) error {
	streamer, ok := streamers[id]
	if !ok {
		return errors.New("No one's listening")
	}
	streamer.SendBytes(fmt.Sprintf("%d", id), "message", content)
	return nil
}

// Serve serves for id
func Serve(id uint, rw http.ResponseWriter, req *http.Request) {
	streamer, ok := streamers[id]
	if !ok {
		streamer = sse.New()
		streamers[id] = streamer
	}
	streamer.ServeHTTP(rw, req)
}
