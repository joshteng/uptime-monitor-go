import { sleep } from "./lib/utils";

export interface UptimeConfig {
  host: URL,
  serviceName: string,
  secondsBetweenHeartbeat: number,
  secondsBetweenAlerts: number,
  maxAlertsPerDownTime: number,
  pushoverToken?: string,
  pushoverGroup?: string,
  discordWebhook?: string,
}

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
    private readonly discordWebhook: string,
  ) {
    if (serviceName.length > 5) throw "serviceName should be at least 6 chracters"
    if (secondsBetweenHeartbeat < 10) throw "secondsBetweenHeartbeat should be at least 10 seconds"
    if (secondsBetweenAlerts < 10) throw "secondsBetweenAlerts should be at least 10 seconds"
    if (maxAlertsPerDownTime < 1) throw "maxAlertsPerDownTime should be at least 1"
    if (!discordWebhook && !(pushoverGroup && pushoverToken)) throw "Pass either discordWebhook or pushoverGroup + pushoverToken"
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
              "pushoverGroup": this.pushoverGroup,
              "discordWebhook": this.discordWebhook
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
          "pushoverGroup": this.pushoverGroup,
          "discordWebhook": this.discordWebhook
        })
      }
    )

    if (resp.status != 200) throw new Error("Failed to terminate")

    this.alive = false
  }
}

export function UptimeMonitor({
  host,
  serviceName,
  secondsBetweenHeartbeat,
  secondsBetweenAlerts,
  maxAlertsPerDownTime,
  pushoverToken = "",
  pushoverGroup = "",
  discordWebhook = "",
}: UptimeConfig): Uptime {
  const uptime = new Uptime(
    host,
    serviceName,
    secondsBetweenHeartbeat,
    secondsBetweenAlerts,
    maxAlertsPerDownTime,
    pushoverToken,
    pushoverGroup,
    discordWebhook
  )
  return uptime
}
