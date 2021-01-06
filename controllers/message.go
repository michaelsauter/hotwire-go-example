package controllers

import (
	"bytes"
	"net/http"

	"github.com/while1malloc0/hotwire-go-example/models"
	"github.com/while1malloc0/hotwire-go-example/pkg/streamer"
)

// MessagesController implements Controller functionality for the Message model
type MessagesController struct {
}

// New renders a form for creating a new Message
func (*MessagesController) New(w http.ResponseWriter, r *http.Request) {
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	responseData := map[string]interface{}{"Room": room}
	render.HTML(w, http.StatusOK, "messages/new", responseData)
}

// Create creates a new Message
func (*MessagesController) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	message := &models.Message{
		Content: r.FormValue("message[content]"),
		Room:    *room,
	}

	err = models.CreateMessage(message)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "text/html; turbo-stream; charset=utf-8")
	var content bytes.Buffer
	responseData := map[string]interface{}{"Message": message}
	err = render.HTML(&content, http.StatusCreated, "messages/create.turbostream", responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	streamer.Publish(room.ID, content.Bytes())
}

// SSE opens a persistent WebSocket connection and subscribes it for updates to its room
func (*MessagesController) SSE(w http.ResponseWriter, r *http.Request) {
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	streamer.Serve(room.ID, w, r)
}
