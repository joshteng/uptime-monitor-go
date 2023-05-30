package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	ServiceName             string `gorm:"uniqueIndex"`
	SecondsBetweenHeartbeat int64
	SecondsBetweenAlerts    int64
	MaxAlertsPerDownTime    int64
	NumberOfAlerts          int64
	IsActive                bool
	PushoverToken           string
	PushoverGroup           string
	LastAlertAt             time.Time
}

func CreateOrUpdateRecord(record Record) Record {
	var _record Record

	if err := DB.Where("service_name = ?", record.ServiceName).First(&_record).Error; err != nil {
		log.Printf("Record not found: %v", err)

		DB.Create(&record)
		return record
	}

	DB.Model(&_record).UpdateColumns(map[string]interface{}{
		"service_name":              record.ServiceName,
		"updated_at":                time.Now(),
		"seconds_between_heartbeat": record.SecondsBetweenHeartbeat,
		"seconds_between_alerts":    record.SecondsBetweenAlerts,
		"max_alerts_per_down_time":  record.MaxAlertsPerDownTime,
		"number_of_alerts":          record.NumberOfAlerts,
		"is_active":                 record.IsActive,
		"last_alert_at":             record.LastAlertAt,
		"pushover_token":            record.PushoverToken,
		"pushover_group":            record.PushoverGroup,
	})
	return _record
}

func UpdateRecord(record Record) Record {
	DB.Model(&record).
		Omit("updated_at").
		UpdateColumns(map[string]interface{}{
			"number_of_alerts": record.NumberOfAlerts,
			"is_active":        record.IsActive,
			"last_alert_at":    record.LastAlertAt})

		// Using struct to update will not set IsActive to false. Bizarre but it's a design
		// UpdateColumns(Record{
		// 	NumberOfAlerts: record.NumberOfAlerts,
		// 	IsActive:       false,
		// 	LastAlertAt:    record.LastAlertAt,
		// })

	return record
}
