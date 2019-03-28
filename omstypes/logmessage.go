package omstypes

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	omsgo "github.com/dtzar/oms-go/oms_data_collector"
)

// LogMessage struct is used to construct a message to be sent to OMS
type LogMessage struct {
	Timestamp             time.Time `json:"timestamp"`
	RequestURI            string    `json:"requestURI"`
	Level                 string    `json:"level"`
	Stage                 string    `json:"stage"`
	Verb                  string    `json:"verb"`
	UserAgent             string    `json:"userAgent"`
	SourceIPAddress       string    `json:"sourceIPAddress"`
	FullUserName          string    `json:"fullUserName"`
	UserDisplayName       string    `json:"userDisplayName"`
	UserPrincipalName     string    `json:"userPrincipalName"`
	UserGroups            string    `json:"userGroups"`
	ResourceType          string    `json:"resourceType"`
	ResourceName          string    `json:"resourceName"`
	ResourceNamespace     string    `json:"resourceNamespace"`
	ResponseStatus        string    `json:"responseStatus"`
	ResponseReason        string    `json:"responseReason"`
	ResponseCode          int32     `json:"responseCode"`
	AuthorizationDecision string    `json:"authorizationDecision"`
	AuthorizationReason   string    `json:"authorizationReason"`
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

// Print outputs all the values of the LogMessage object to the STDOUT stream
func (lm *LogMessage) Print() {
	// Prefix the log entry with "[AuditRecord]" so we can distinguish between real data and error messages in the logs
	log.Println(
		fmt.Sprintf(
			"[AuditRecord] %s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%d,%s,%s",
			lm.Timestamp,
			lm.RequestURI,
			lm.Level,
			lm.Stage,
			lm.Verb,
			lm.UserAgent,
			lm.SourceIPAddress,
			lm.FullUserName,
			lm.UserDisplayName,
			lm.UserPrincipalName,
			lm.UserAgent,
			lm.ResourceType,
			lm.ResourceName,
			lm.ResourceNamespace,
			lm.ResponseStatus,
			lm.ResponseReason,
			lm.ResponseCode,
			lm.AuthorizationDecision,
			lm.AuthorizationReason))
}
