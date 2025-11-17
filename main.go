package main

import (
  "fmt"
  "power/configuration"
  "power/server"
  "power/enphase"
  "net/http"
)

func main() {
  token:=configuration.Config.Token
  if (token == "") {
    fmt.Printf("getting token")
    token = enphase.GetEnphaseToken()
  }

  mux := http.NewServeMux()
  server.SetupRoutes(mux)
  errs := make(chan error)
  go enphase.StartMetricDataStream(token, errs);
  go listenHTTP(mux, errs)
  fmt.Printf("%v",<-errs)
}


func listenHTTP(mux *http.ServeMux, errs chan<- error) {
  errs <- http.ListenAndServe(fmt.Sprintf(":%v",configuration.Config.HttpPort), mux)
}
