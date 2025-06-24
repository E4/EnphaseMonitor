package server

import (
  "fmt"
  "time"
  "power/filesystem"
  "net/http"
  "github.com/gorilla/websocket"
  "path/filepath"
  "power/enphase"
  "os"
  "encoding/json"
)

var upgrader = websocket.Upgrader {
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}

func SetupRoutes(mux *http.ServeMux) {
  /* register endpoints */
  mux.HandleFunc("/time", endpointTime)
  mux.HandleFunc("/meter", endpointMeter)
  mux.HandleFunc("/error", endpointError)
  mux.HandleFunc("/stream", websocketUpgrade)
  mux.Handle("/", staticFileServer("./static"))
}


/* simple endpoint to display time */
func endpointTime(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "%s", time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
}


func endpointMeter(w http.ResponseWriter, r *http.Request) {

  data, err := json.Marshal(enphase.LatestMeterData)
  if err != nil {
    http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(data)
}


/* */
func endpointError(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "Bad Request", 400)
}


/* serve static files (development only) */
func staticFileServer(servePath string) http.Handler {
  fs := http.FileServer(filesystem.RestrictedFileSystem{http.Dir(servePath), true})
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      fullPath := filepath.Join(servePath, r.URL.Path)
    info, err := os.Stat(fullPath)
    if err != nil || info.IsDir() {
      http.ServeFile(w, r, filepath.Join(servePath, "index.html"))
      return
    }
    fs.ServeHTTP(w, r)
  })
}


/* websockets endpoint */
func websocketUpgrade(w http.ResponseWriter, r *http.Request) {
  upgrader.CheckOrigin = func(r *http.Request) bool { return true }
  connection, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    return
  }
  enphase.StreamToConnection(connection);
}

