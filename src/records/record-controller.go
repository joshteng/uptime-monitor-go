package records

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"uptime-monitor/src/models"
	"uptime-monitor/src/utils"
)

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response := map[string]interface{}{
			"success": false,
			"message": "No routes found",
		}

		utils.ReturnJsonResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var recordInput RecordInputDto
	if err := json.Unmarshal(reqBody, &recordInput); err != nil {
		response := map[string]interface{}{
			"success": false,
			"message": "Error parsing request body",
		}

		utils.ReturnJsonResponse(w, http.StatusBadRequest, response)
		return
	}

	if err := validateRecordInput(recordInput); err != nil {
		errorMessages := make([]string, 0)

		for i := 0; i < len(err); i++ {
			errorMessages = append(errorMessages, err[i].Error())
		}

		response := map[string]interface{}{
			"success": false,
			"message": strings.ReplaceAll(fmt.Sprintf("%+q", errorMessages), "\" \"", "\",\""),
		}
		utils.ReturnJsonResponse(w, http.StatusBadRequest, response)
		return
	}

	record := models.Record{
		ServiceName:             recordInput.ServiceName,
		SecondsBetweenHeartbeat: recordInput.SecondsBetweenHeartbeat,
		SecondsBetweenAlerts:    recordInput.SecondsBetweenAlerts,
		MaxAlertsPerDownTime:    recordInput.MaxAlertsPerDownTime,
		NumberOfAlerts:          0,
		IsActive:                true,
		LastAlertAt:             time.Now().Add(-time.Second * time.Duration(recordInput.SecondsBetweenAlerts)),
		PushoverToken:           recordInput.PushoverToken,
		PushoverGroup:           recordInput.PushoverGroup,
		DiscordWebhook:          recordInput.DiscordWebhook,
	}

	updatedRecord := models.CreateOrUpdateRecord(record)
	respBody := transformToDto(updatedRecord)
	recordJSON, err := json.Marshal(&respBody)
	if err != nil {

		response := map[string]interface{}{
			"success": false,
			"message": "Error parsing record",
		}
		utils.ReturnJsonResponse(w, http.StatusInternalServerError, response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(recordJSON)
}

func validateRecordInput(recordInput RecordInputDto) []error {
	var _errors []error
	if len(recordInput.ServiceName) == 0 {
		_errors = append(_errors, errors.New("serviceName cannot be blank"))
	}

	if recordInput.SecondsBetweenHeartbeat == 0 {
		_errors = append(_errors, errors.New("secondsBetweenHeartbeat cannot be 0"))
	}

	// if recordInput.MaxAlertsPerDownTime == 0 {
	// 	_errors = append(_errors, errors.New("maxAlertsPerDownTime cannot be 0"))
	// }

	if recordInput.SecondsBetweenAlerts == 0 {
		_errors = append(_errors, errors.New("secondsBetweenAlerts cannot be 0"))
	}

	if len(_errors) > 0 {
		return _errors
	} else {
		return nil
	}
}

func transformToDto(record models.Record) RecordInputDto {
	return RecordInputDto{
		ServiceName:             record.ServiceName,
		SecondsBetweenHeartbeat: record.SecondsBetweenHeartbeat,
		MaxAlertsPerDownTime:    record.MaxAlertsPerDownTime,
		SecondsBetweenAlerts:    record.SecondsBetweenAlerts,
	}
}
