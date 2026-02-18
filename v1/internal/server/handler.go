package server

import (
	"embed"
	"io/fs"
	"net/http"

	"nhooyr.io/websocket"
)

//go:embed web
var webFS embed.FS

// webSubFS is the embedded web/ directory rooted at "/", so files are
// accessible as "/index.html", "/dashboard.js", etc.
var webSubFS = func() fs.FS {
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		// Programming error: the embedded directory must always exist.
		panic("server: failed to sub embedded web FS: " + err.Error())
	}
	return sub
}()

// NewHandler returns an HTTP handler that:
//   - upgrades /ws connections to WebSocket and delegates them to the hub
//   - serves embedded static files (index.html, dashboard.js, dashboard.css) at /
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
	mux.Handle("/", http.FileServer(http.FS(webSubFS)))
	return mux
}
