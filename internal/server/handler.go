package server

import (
	"net/http"

	"nhooyr.io/websocket"
)

// NewHandler returns an HTTP handler that upgrades connections to WebSocket
// and delegates them to the hub.
func NewHandler(hub *Hub) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true, // Allow connections from any origin.
		})
		if err != nil {
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "")
		hub.ServeClient(r.Context(), conn)
	})
	return mux
}
