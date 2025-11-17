
# Settings

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

# Starting

Install dependencies using Go `go get`, and start with `go run main.go` and navigate to http://localhost:[httpport]/ with a local browser.


# URLs

`/meter` Will show you the latest data that was received. It's a JSON formed data packet that provides current values for production, net-consumption and total-consumption. The data is very detailed, and per phase.

`/stream`  This is a websocket endpoint that streams out the data as it comes in.

`/index.html`  This program also has an HTTP server that will serve contents from the /static folder. There's a program included that attaches to the stream and shows the data in a neat chart.

# Screenshot

![Screenshot](/docs/screenshot.png?raw=true)