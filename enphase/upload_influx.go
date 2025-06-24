package enphase

import (
  "bytes"
  "fmt"
  "net/http"
  "time"
)


func uploadToInfluxDB(m MeterData, influxURL, org, bucket, token string) error {
  lines := []string{}
  timestamp := time.Now().UnixNano()

  // Only include phase C if voltage is non-zero in both sections
  includePhaseC := !(m.Production.PHC.V == 0 || m.TotalConsumption.PHC.V == 0)

  // Define the two data sources to upload
  entries := []struct {
    prefix  string
    section ThreePhaseSection
  }{
    {"production_phase", m.Production},
    {"consumption_phase", m.TotalConsumption},
  }

  for _, entry := range entries {
    phaseMap := map[string]PhaseData{
      "a": entry.section.PHA,
      "b": entry.section.PHB,
    }

    if includePhaseC {
      phaseMap["c"] = entry.section.PHC
    }

    for phase, data := range phaseMap {
      src := fmt.Sprintf("%s_%s", entry.prefix, phase)

      line := fmt.Sprintf(
        "meter_data,src=%s active_power=%f,reactive_power=%f,apparent_power=%f,voltage=%f,current=%f,power_factor=%f,frequency=%f %d",
        src,
        data.P, data.Q, data.S, data.V, data.I, data.PF, data.F,
        timestamp,
      )
      lines = append(lines, line)
    }
  }

  if len(lines) == 0 {
    return fmt.Errorf("no data to upload (all values filtered out)")
  }

  payload := []byte(lines[0] + "\n")
  for i := 1; i < len(lines); i++ {
    payload = append(payload, []byte(lines[i]+"\n")...)
  }

  req, err := http.NewRequest("POST", influxURL, bytes.NewBuffer(payload))
  if err != nil {
    return fmt.Errorf("creating request: %w", err)
  }

  req.Header.Set("Authorization", "Token "+token)
  req.Header.Set("Content-Type", "text/plain")
  q := req.URL.Query()
  q.Add("org", org)
  q.Add("bucket", bucket)
  q.Add("precision", "ns")
  req.URL.RawQuery = q.Encode()

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return fmt.Errorf("sending request: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusNoContent {
    return fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, resp.Status)
  }

  return nil
}