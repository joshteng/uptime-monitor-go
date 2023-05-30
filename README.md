# Uptime Monitor

A simple downtime monitoring and alert service for all your processes!

Great for processes and microservices that do not expose any ports or endpoints for monitoring via fabulous services such as Uptime Robot (free version). Think push instead of poll.

This is a rewrite of https://github.com/joshteng/uptime-monitor in Golang.

## How it works
It sends a [Pushover](https://pushover.net/) notification whenever your service stops sending a pulse (see endpoint below).

This app runs a loop to to check which services are down in addition to a HTTP server that accepts updates from services.

There is a single endpoint application that needs no configuration or set up for whatever new service or process you want to monitor. Just send a HTTP request (see endpoint below) and it will do it's magic.

## Environment Variables
```
SQLITE_PATH
```

## Deployment
This project was written in Golang. You can host it like how you host any Go projects.

Quick start:
```sh
go mod download
go build -o uptime-monitor
./uptime-monitor
```

By default, it runs on port 8080.

Alternatively, a `Dockerfile` and a `docker-compose.example.yml` is included.

Also, an example `fly.example.toml` is provided for deployment to fly.io

## Client Side
Right now, there is only a JS SDK https://www.npmjs.com/package/uptime-monitor-sdk source code in this repo: `/packages/sdk`

Check out the README on NPMJS on how to use it.

Otherwise, see below:

## Creating or Updating a service (to let our Go app know that our service exists or is still alive)
Send a HTTP POST request to /records
```curl
curl -X POST -H "Content-Type: application/json" \
  -d '{"serviceName": "<YOUR_SERVICE_NAME>", "secondsBetweenHeartbeat": 60, "secondsBetweenAlerts": 600, "maxAlertsPerDownTime": 10, "pushoverToken": "<YOUR_PUSHOVER_TOKEN>, "pushoverGroup": "<YOUR_PUSHOVER_GROUP>"}' \
  http://localhost:3000/records
```

## Endpoint Variables Defined
|Parameters|Type|Description|
|---|---|---|
|serviceName|string|Name of your app or service or script (needs to be unique)|
|secondsBetweenHeartbeat|number|A whole number to specify how often you want to monitor your service and correspondingly ping the backend. If the server stops receiving after this duration, it will start sending out Push Notification alerts via Pushover|
|secondsBetweenAlerts|number|A whole number of how often (in seconds) you want to receive Push Notification alerts|
|maxAlertsPerDownTime|number|A whole number of the max number of alerts you will receive before notification stops|
|pushoverToken|string|Obtain this from pushover.net|
|pushoverGroup|string|Obtain this from pushover.net|

## Security
Feel free to implement your own basic auth or whatever. I operate this in a controlled private network environment.

## Reliability
Won't it be bad if your uptime monitor goes down instead of your app?

Typical production deployment practices will help prevent that from happening. Or you could deploy a new instance of this app for every instance of this app you deploy. ;p
