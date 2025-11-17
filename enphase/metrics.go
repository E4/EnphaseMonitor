package enphase

import (
  "bufio"
  "encoding/json"
  "crypto/tls"
  "fmt"
  "io"
  "net/http"
  "net/http/cookiejar"
  "strings"
  "power/configuration"
  "time"
)


// PhaseData holds p, q, s, v, i, pf, f for a single phase.
type PhaseData struct {
  P  float64 `json:"p"`
  Q  float64 `json:"q"`
  S  float64 `json:"s"`
  V  float64 `json:"v"`
  I  float64 `json:"i"`
  PF float64 `json:"pf"`
  F  float64 `json:"f"`
}

// ThreePhaseSection holds data for ph-a, ph-b, ph-c.
type ThreePhaseSection struct {
  PHA PhaseData `json:"ph-a"`
  PHB PhaseData `json:"ph-b"`
  PHC PhaseData `json:"ph-c"`
}

// MeterData is the full structure of each streamed data packet.
type MeterData struct {
  Production        ThreePhaseSection `json:"production"`
  NetConsumption    ThreePhaseSection `json:"net-consumption"`
  TotalConsumption  ThreePhaseSection `json:"total-consumption"`
}

var LatestMeterData MeterData;

func StartMetricDataStream(token string, errs chan<- error) {
  jar, _ := cookiejar.New(nil)
  client := &http.Client{Jar: jar}

  authURL := "https://"+configuration.Config.EnvoyAddress+"/auth/check_jwt"
  streamURL := "https://"+configuration.Config.EnvoyAddress+"/stream/meter"
  client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return nil } // prevent auto redirect
  insecureClient := *client
  insecureClient.Transport = &http.Transport{TLSClientConfig: insecureTLS()}
  _ = errs

  for {
    if err := authorize(&insecureClient, authURL, token); err != nil {
      fmt.Printf("Auth check failed: %v. Retrying in 10s.\n", err)
      time.Sleep(10 * time.Second)
      continue
    }

    streamReq, err := http.NewRequest("GET", streamURL, nil)
    if err != nil {
      fmt.Printf("Failed to create stream request: %v. Retrying in 5s.\n", err)
      time.Sleep(5 * time.Second)
      continue
    }

    streamResp, err := insecureClient.Do(streamReq)
    if err != nil {
      fmt.Printf("Stream request failed: %v. Retrying in 10s.\n", err)
      time.Sleep(10 * time.Second)
      continue
    }

    if err := streamAndParse(streamResp.Body); err != nil {
      fmt.Printf("Stream disconnected: %v. Reconnecting...\n", err)
    } else {
      fmt.Printf("Stream ended without error. Reconnecting...\n")
    }

    time.Sleep(5 * time.Second)
  }
}


func streamAndParse(streamRespBody io.ReadCloser) error {
  defer streamRespBody.Close()
  scanner := bufio.NewScanner(streamRespBody)
  for scanner.Scan() {
    line := scanner.Text()
    if !strings.HasPrefix(line, "data: ") {
      continue
    }
    jsonStr := strings.TrimPrefix(line, "data: ")
    if err := json.Unmarshal([]byte(jsonStr), &LatestMeterData); err != nil {
      fmt.Printf("Failed to parse JSON: %v\n", err)
      continue
    }
    // fmt.Printf("%+v\n", LatestMeterData)
    uploadData(LatestMeterData)
  }
  if err := scanner.Err(); err != nil {
    return fmt.Errorf("scanner finished %w", err)
  }
  return fmt.Errorf("stream closed")
}


func uploadData(data MeterData) {
  if configuration.Config.InfluxURL!="" {
    err := uploadToInfluxDB(data,
      configuration.Config.InfluxURL,
      configuration.Config.InfluxOrgId,
      configuration.Config.InfluxBucket,
      configuration.Config.InfluxToken,
    )
    if err != nil {
      fmt.Printf("Upload failed: %v\n", err)
    }
  }
  uploadToWebsocket(data);
}


func extractTextareaContent(html string) string {
  start := strings.Index(html, "<textarea")
  if start == -1 {
    return ""
  }
  start = strings.Index(html[start:], ">")
  if start == -1 {
    return ""
  }
  start += strings.Index(html, "<textarea") + 1
  end := strings.Index(html[start:], "</textarea>")
  if end == -1 {
    return ""
  }

  return strings.TrimSpace(html[start : start+end])
}


func authorize(client *http.Client, url string, token string) error {
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return err
  }
  req.Header.Set("Authorization", "Bearer "+token)
  _, err = client.Do(req)
  return err
}


func insecureTLS() *tls.Config {
  return &tls.Config{InsecureSkipVerify: true}
}
