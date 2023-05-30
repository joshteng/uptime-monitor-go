package records

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"uptime-monitor/src/models"
)

func Monitor() {
	for {
		detect()
		time.Sleep(5 * time.Second)
	}
}

func detect() {
	now := time.Now()

	var records []models.Record

	if err := models.DB.Find(&records, "is_active = ?", true).Error; err != nil {
		log.Fatalf("Retrieving records failed %v", err)
	}

	for i := 0; i < len(records); i++ {
		record := records[i]

		secondsFromLastUpdate := now.Sub(record.UpdatedAt).Seconds()
		secondsFromLastAlert := now.Sub(record.LastAlertAt).Seconds()

		if secondsFromLastUpdate > float64(record.SecondsBetweenHeartbeat) && secondsFromLastAlert > float64(record.SecondsBetweenAlerts) {
			msg := record.ServiceName + " has been down for " + fmt.Sprintf("%f", secondsFromLastUpdate) + " seconds! " + fmt.Sprintf("%f", secondsFromLastUpdate/60) + " minutes"
			log.Println(msg)

			// send notification
			sendAlert(msg, record.PushoverToken, record.PushoverGroup)

			record.LastAlertAt = now

			if record.NumberOfAlerts+1 >= record.MaxAlertsPerDownTime {
				record.IsActive = false
			} else {
				record.NumberOfAlerts = record.NumberOfAlerts + 1
			}

			models.UpdateRecord(record)
		}

	}
}

type PushoverMessage struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Message string `json:"message"`
}

func sendAlert(msg, pushoverToken, pushoverGroup string) {
	pushoverMessage := PushoverMessage{
		Token:   pushoverToken,
		User:    pushoverGroup,
		Message: msg,
	}

	reqBody, err := json.Marshal(pushoverMessage)
	bodyReader := bytes.NewReader(reqBody)

	if err != nil {
		log.Fatal("Failed to create request body")
	}

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, "https://api.pushover.net/1/messages.json", bodyReader)
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal(err)
	}
}
