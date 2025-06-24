package enphase

import (
  "fmt"
  "io"
  "net/http"
  "net/http/cookiejar"
  "net/url"
  "power/configuration"
)

func GetEnphaseToken() string {
  // Create a cookie jar to store cookies across requests
  jar, _ := cookiejar.New(nil)
  client := &http.Client{Jar: jar}

  // Step 1: POST to login to get session cookies
  loginURL := "https://entrez.enphaseenergy.com/login"
  loginData := url.Values{
    "username":       {configuration.Config.UserName},
    "password":       {configuration.Config.Password},
    "codeChallenge":  {""},
    "redirectUri":    {""},
    "client":         {""},
    "clientId":       {""},
    "authFlow":       {"entrezSession"},
    "serialNum":      {""},
    "grantType":      {""},
    "state":          {""},
    "invalidSerialNum": {""},
  }

  _, err := client.PostForm(loginURL, loginData)
  if err != nil {
    panic(fmt.Sprintf("Login request failed: %v", err))
  }

  // Step 2: Get the token from the textarea response
  tokenURL := "https://entrez.enphaseenergy.com/entrez_tokens"
  tokenData := url.Values{
    "Site":      {configuration.Config.Site},
    "serialNum": {configuration.Config.SerialNum},
  }
  resp, err := client.PostForm(tokenURL, tokenData)
  if err != nil {
    panic(fmt.Sprintf("Token request failed: %v", err))
  }
  defer resp.Body.Close()

  // Extract token from <textarea>
  bodyBytes, err := io.ReadAll(resp.Body)
  if err != nil {
    panic(fmt.Sprintf("Failed reading token body: %v", err))
  }
  token := extractTextareaContent(string(bodyBytes))
  if token == "" {
    panic("Failed to extract token from textarea")
  }
  return token
}
