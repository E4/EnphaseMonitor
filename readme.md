Create a file `config/system.ini` with the desired configuration. It needs to have the following values:

```
localip=[the local ip address of the enphase box]
httpport=[the port this program will use to serve http]
```

Either provide username, password, site and serial number for automatic token generation, or provide token:
```
username=[username for enphase - used for generating token]
password=[password for enphase - used for generating token]
site=[customer's name - used for generating token]
serialnum=[from enphase web site - used for generating token]
```

Bypass token generation by providing token directly
```
token=[the token generated from Enphase web site]
```

Optional for uploading the data to InfluxDB
```
influx_token=
influx_orgid=
influx_bucket=
influx_url=
```


Install dependencies using Go, and start with `go run main.go` and navigate to http://localhost:[httpport]/ with a local browser.
