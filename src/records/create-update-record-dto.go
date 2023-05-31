package records

type RecordInputDto struct {
	ServiceName             string `json:"serviceName"`
	SecondsBetweenHeartbeat int64  `json:"secondsBetweenHeartbeat"`
	MaxAlertsPerDownTime    int64  `json:"maxAlertsPerDownTime"`
	SecondsBetweenAlerts    int64  `json:"secondsBetweenAlerts"`
	PushoverToken           string `json:"pushoverToken"`
	PushoverGroup           string `json:"pushoverGroup"`
	DiscordWebhook          string `json:"discordWebhook"`
}
