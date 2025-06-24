package configuration

import (
  "io/ioutil"
  "strings"
  "strconv"
)

type ConfigurationStruct struct {
  UserName string
  Password string
  Site string
  SerialNum string
  EnvoyAddress string
  HttpPort uint16
  Token string
  InfluxURL string
  InfluxToken string
  InfluxOrgId string
  InfluxBucket string
}

func LoadConfigurationFromFile(path string) *ConfigurationStruct{
  dat, err := ioutil.ReadFile(path)
  if err!=nil { dat, err = ioutil.ReadFile("../" + path) }
  if err!=nil { panic(err) }
  return parse(string(dat))
}

func parse(data string) *ConfigurationStruct {
  var config *ConfigurationStruct = new(ConfigurationStruct)
  var name,value string
  z := strings.Split(data,"\n")
  for i:=0;i<len(z);i++ {
    dz := strings.SplitN(z[i],"=",2)
    name = strings.TrimSpace(dz[0])
    if len(name)==0 { continue }
    if len(dz)==1 { value = "" } else { value = strings.TrimSpace(dz[1]) }
    setConfigurationValue(config, name, value)
  }
  return config
}

func setConfigurationValue(config *ConfigurationStruct, name string, value string) {
  switch (strings.ToLower(name)) {
    case "username":
      config.UserName = value
    case "password":
      config.Password = value
    case "site":
      config.Site = value
    case "serialnum":
      config.SerialNum = value
    case "localip":
      config.EnvoyAddress = value
    case "httpport":
      config.HttpPort = parseUint16(value, 10)
    case "token":
      config.Token = value
    case "influx_token":
      config.InfluxToken = value
    case "influx_orgid":
      config.InfluxOrgId = value
    case "influx_bucket":
      config.InfluxBucket = value
    case "influx_url":
      config.InfluxURL = value
  }
}

func parseUint16(str string, base int) uint16 {
  i,_ := strconv.ParseUint(str, base, 16)
  return uint16(i)
}

var Config *ConfigurationStruct;

func init() {
  Config = LoadConfigurationFromFile("./config/system.ini")
}

