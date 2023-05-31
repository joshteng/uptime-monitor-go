import { sleep } from "./lib/utils";

class FetchError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'FetchError';
  }
}

class Uptime {
  private alive: boolean

  constructor(
    private readonly host: URL,
    private readonly serviceName: string,
    private readonly secondsBetweenHeartbeat: number,
    private readonly secondsBetweenAlerts: number,
    private readonly maxAlertsPerDownTime: number,
    private readonly pushoverToken: string,
    private readonly pushoverGroup: string,
  ) {
    if (secondsBetweenHeartbeat < 15) throw "secondsBetweenHeartbeat should be at least 15 seconds"
    if (secondsBetweenAlerts < 60) throw "secondsBetweenAlerts should be at least 60 seconds"
    if (maxAlertsPerDownTime < 1) throw "maxAlertsPerDownTime should be at least 1"
    if (!pushoverGroup || !pushoverToken) throw "Missing pushover token and / or group"
    this.alive = true;
  }

  public init() {
    fetch(this.host, {
      headers: {
        'Content-Type': 'application/json'
      }
    }).then(resp => resp.json()).then(res => {
      if (!res.success) throw new FetchError("Uptime Monitor: Invalid HOST or application not running")
    }).catch(err => { throw new FetchError("Uptime Monitor: Invalid HOST or application not running") })

    this.heartbeat()

    return this
  }

  private async heartbeat() {
    while (this.alive) {
      try {
        await fetch(
          this.host + "records",
          {
            method: 'POST',
            headers: {
              'Content-type': 'application/json'
            },
            body: JSON.stringify({
              "serviceName": this.serviceName,
              "secondsBetweenHeartbeat": this.secondsBetweenHeartbeat,
              "secondsBetweenAlerts": this.secondsBetweenAlerts,
              "maxAlertsPerDownTime": this.maxAlertsPerDownTime,
              "pushoverToken": this.pushoverToken,
              "pushoverGroup": this.pushoverGroup
            })
          }
        )

        await sleep((this.secondsBetweenHeartbeat - Math.floor(this.secondsBetweenHeartbeat / 2)) * 1000)
      } catch (err) {
        console.log("Heartbeat failed. Retrying.....")
        console.debug(err)
      }
    }
  }

  public async terminate() {
    const resp = await fetch(
      this.host + "records",
      {
        method: 'POST',
        headers: {
          'Content-type': 'application/json'
        },
        body: JSON.stringify({
          "serviceName": this.serviceName,
          "secondsBetweenHeartbeat": this.secondsBetweenHeartbeat,
          "secondsBetweenAlerts": this.secondsBetweenAlerts,
          "maxAlertsPerDownTime": 0,
          "pushoverToken": this.pushoverToken,
          "pushoverGroup": this.pushoverGroup
        })
      }
    )

    if (resp.status != 200) throw new Error("Failed to terminate")

    this.alive = false
  }
}

export function UptimeMonitor(
  host: URL,
  serviceName: string,
  secondsBetweenHeartbeat: number,
  secondsBetweenAlerts: number,
  maxAlertsPerDownTime: number,
  pushoverToken: string,
  pushoverGroup: string,
): Uptime {
  const uptime = new Uptime(
    host,
    serviceName,
    secondsBetweenHeartbeat,
    secondsBetweenAlerts,
    maxAlertsPerDownTime,
    pushoverToken,
    pushoverGroup
  )
  return uptime
}
