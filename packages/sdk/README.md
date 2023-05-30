# Uptime Monitor JS SDK

This SDK is used on the client side to repeatedly ping the Uptime Monitor server. When the client side terminates / crashes and the server stops receiving pings, alerts will be pushed via Pushover.

Pushover is configured on the server side.

You will need to host the Uptime Monitor backend and point a domain to it.

## Usage
```bash
npm i uptime-monitor-sdk
```

## To maintain a session
```js
import { UptimeMonitor } from "uptime-monitor-sdk";

const url = new URL("https://uptime.example.com")

UptimeMonitor(url, 'Name of App', 10, 60, 10, "<PUSHOVER_TOKEN>", "<PUSHOVER_GROUP>").init() // notice the `init()`!
```

## To terminate a session
```js
import { UptimeMonitor } from "uptime-monitor-sdk";

const url = new URL("https://uptime.example.com")

const uptime = UptimeMonitor(url, 'Name of App', 10, 60, 10, "<PUSHOVER_TOKEN>", "<PUSHOVER_GROUP>")

uptime.terminate()
```

## Full Example
```js
import { UptimeMonitor } from "uptime-monitor-sdk";

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

const url = new URL("https://uptime.example.com")

const uptime = UptimeMonitor(
  url,
  'Example',
  10,
  60,
  10,
  "<PUSHOVER_TOKEN>",
  "<PUSHOVER_GROUP>",
).init()

async function main(fn) {
  while (true) {
    console.log("Still alive....")
    const rand = Math.floor(Math.random() * 1000);

    if (rand > 900) break;

    await sleep(10 * 1000)
  }

  fn()
}

main(() => uptime.terminate())
```

## Constructor Parameters
```js
UptimeMonitor(url, serviceName, secondsBetweenHeartbeat, secondsBetweenAlerts, maxAlertsPerDownTime, pushoverToken, pushoverGroup)
```

|Parameters|Type|Description|
|---|---|---|
|url|URL|The backend URL for Uptime Monitor|
|serviceName|string|Name of your app or service or script (needs to be unique)|
|secondsBetweenHeartbeat|number|A whole number to specify how often you want to monitor your service and correspondingly ping the backend. If the server stops receiving after this duration, it will start sending out Push Notification alerts via Pushover|
|secondsBetweenAlerts|number|A whole number of how often (in seconds) you want to receive Push Notification alerts|
|maxAlertsPerDownTime|number|A whole number of the max number of alerts you will receive before notification stops|
|pushoverToken|string|Obtain this from pushover.net|
|pushoverGroup|string|Obtain this from pushover.net|

## Calling `init`

You need to call `init()` in order to do periodic pings to uptime backend!

## Calling `terminate`

You should call `terminate` when the process ends or is no longer needed. This prevents you from getting unnecessary notifications.
