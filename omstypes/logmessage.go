package omstypes

import (
	"encoding/json"
	"time"

	omsgo "github.com/dtzar/oms-go/oms_data_collector"
)

// LogMessage struct is used to construct a message to be sent to OMS
type LogMessage struct {
	Timestamp         time.Time `json:"timestamp"`
	RequestURI        string    `json:"requestURI"`
	Verb              string    `json:"verb"`
	UserDisplayName   string    `json:"userDisplayName"`
	UserPrincipalName string    `json:"userPrincipalName"`
	ResourceType      string    `json:"resourceType"`
	ResourceName      string    `json:"resourceName"`
}

// PostToOMS attempts to actually post the message to OMS for the given log type
func (lm *LogMessage) PostToOMS(omsLogClient omsgo.OmsLogClient, omsLogType string) error {
	buffer, err := json.Marshal(lm)
	if err != nil {
		return err
	}

	responseErr := omsLogClient.PostData(&buffer, omsLogType)
	if responseErr != nil {
		return responseErr
	}

	return nil
}
