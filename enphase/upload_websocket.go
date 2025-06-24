package enphase

import (
  "github.com/gorilla/websocket"
  "encoding/json"
)

var activeConnection * websocket.Conn;

func StreamToConnection(connection *websocket.Conn) {
  if activeConnection!=nil {
    activeConnection.Close();
  }

  activeConnection = connection;

  for {
    _, _, err := connection.ReadMessage()
    if err != nil {
      connection.Close()
      connection = nil
      activeConnection = nil;
      return
    }
  }
}


func uploadToWebsocket(data MeterData) {
  if activeConnection==nil { return }
  jsonData, err := json.Marshal(data)
  if err != nil { return }
  activeConnection.WriteMessage(websocket.TextMessage, jsonData)
}