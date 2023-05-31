# Uptime Monitor JS SDK

This SDK is used on the client side to repeatedly ping the Uptime Monitor server. When the client side terminates / crashes and the server stops receiving pings, alerts will be pushed via Pushover and/or Discord.

## Usage
```bash
npm i uptime-monitor-sdk
```

## To maintain a session
```js
import { UptimeMonitor, type UptimeConfig } from "uptime-monitor-sdk";

const uptimeConfig: UptimeConfig = {
  host: new URL("https://uptime.example.com"),
  serviceName: 'Example',
  secondsBetweenHeartbeat: 60, // expected heartbeat. if we miss a heartbeat over 60 seconds, an alert is pushed
  secondsBetweenAlerts: 180, // every 3 minutes
  maxAlertsPerDownTime: 10, // maximum number of times you want to be alerted
  pushoverToken: "", // obtain from pushover (optional)
  pushoverGroup: "", // obtain from pushover (optional)
  discordWebhook: "" // obtain from discord (optional)
} // note: you can use either Pushover or Discord or both

UptimeMonitor(uptimeConfig).init() // notice the `init()`!
```

## To terminate a session
```js
// like the above code but we call terminate instead of init()
UptimeMonitor(uptimeConfig).terminate()

// for convenient, init() method returns an instance of the UptimeMonitor
const uptime = UptimeMonitor(uptimeConfig)

// some other code

uptime.terminate()
```

## Full Example
```js
import { UptimeMonitor } from "uptime-monitor-sdk";

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

const uptimeConfig: UptimeConfig = {
  host: new URL("https://uptime.example.com"),
  serviceName: 'Example',
  secondsBetweenHeartbeat: 60, // expected heartbeat. if we miss a heartbeat over 60 seconds, an alert is pushed
  secondsBetweenAlerts: 180, // every 3 minutes
  maxAlertsPerDownTime: 10, // maximum number of times you want to be alerted
  pushoverToken: "", // obtain from pushover (optional)
  pushoverGroup: "", // obtain from pushover (optional)
  discordWebhook: "" // obtain from discord (optional)
} // note: you can use either Pushover or Discord or both

const uptime = UptimeMonitor(uptimeConfig).init()

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

|Parameters|Type|Description|
|---|---|---|
|url|URL|The backend URL for Uptime Monitor|
|serviceName|string|Name of your app or service or script (needs to be unique)|
|secondsBetweenHeartbeat|number|A whole number to specify how often you want to monitor your service and correspondingly ping the backend. If the server stops receiving after this duration, it will start sending out Push Notification alerts via Pushover|
|secondsBetweenAlerts|number|A whole number of how often (in seconds) you want to receive Push Notification alerts|
|maxAlertsPerDownTime|number|A whole number of the max number of alerts you will receive before notification stops|
|pushoverToken|string|Obtain this from pushover.net|
|pushoverGroup|string|Obtain this from pushover.net|
|discordWebhook|URL string|Obtain this from your Discord server|

Note: You can use either Pushover or Discord or both. Simply omit whichever you don't need.

## Calling `init`

You need to call `init()` in order to do periodic pings to uptime backend!

## Calling `terminate`

You should call `terminate` when the process ends or is no longer needed. This prevents you from getting unnecessary notifications.
