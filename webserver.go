package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	graphapi "github.com/carlsoncoder/audit-webhook-go/graphapi"
	kubernetestypes "github.com/carlsoncoder/audit-webhook-go/kubernetestypes"
)

var (
	tenantID            string
	clientID            string
	clientSecret        string
	graphAPIClient      *graphapi.Client
	userTenantURLPrefix string
)

const (
	tenantIDVariableName     = "TENANT_ID"
	clientIDVariableName     = "CLIENT_ID"
	clientSecretVariableName = "CLIENT_SECRET"
)

func kubernetesAuditPostHandler(rw http.ResponseWriter, req *http.Request) {
	now := getCurrentDateTimeString()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(fmt.Sprintf("[%s] ERROR: Unable to read POST body", now))
		log.Println(fmt.Sprintf("%v", err))
		return
	}

	var eventList kubernetestypes.EventList

	err = json.Unmarshal(body, &eventList)
	if err != nil {
		log.Println(fmt.Sprintf("[%s] ERROR: Unable to parse POST to JSON", now))
		log.Println(fmt.Sprintf("%v", err))
		log.Println("Full POST Body:")
		log.Println(string(body[:]))
		return
	}

	// error handling function to capture any errors that occur in the processing in the rest of this function
	defer func() {
		err := recover()
		if err != nil {
			log.Println(fmt.Sprintf("[%s] ERROR: Unable to process kubernetes audit records", now))
			log.Println(fmt.Sprintf("%v", err))
			log.Println("Full POST Body:")
			log.Println(string(body[:]))
		}
	}()

	for _, event := range eventList.Events {
		// We only want to log records that came from an AAD user
		if event.User.IsUserValidAADUser(userTenantURLPrefix) {
			// determine user and user group properties
			userObjectID, userDisplayName, userPrincipalName := event.User.GetUserDetails(userTenantURLPrefix, graphAPIClient)
			userGroups := event.User.GetUserGroups(graphAPIClient)

			// get the source IP address(es) from the event
			sourceIPAddress := event.GetSourceIPAddress()

			// Prefix the log entry with "[AuditRecord]" so we can distinguish between real data and error messages in the logs
			log.Println(
				fmt.Sprintf(
					"[AuditRecord] %s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
					event.StageTimestamp,
					event.Level,
					event.Stage,
					event.RequestURI,
					event.Verb,
					userObjectID,
					userDisplayName,
					userPrincipalName,
					sourceIPAddress,
					userGroups,
					event.ObjectRef.Name,
					event.ObjectRef.Resource,
					event.Annotations.Decision,
					event.Annotations.Reason))
		}
	}
}

func main() {
	// load the necessary parameters from OS environment variables and validate they are present
	loadAndValidateParameters()

	// initialize our graphAPIClient with the parameters
	graphAPIClient = graphapi.NewClient(tenantID, clientID, clientSecret)
	userTenantURLPrefix = fmt.Sprintf("https://sts.windows.net/%s/#", tenantID)

	// setup the handler for POSTing to the /audits endpoint
	http.HandleFunc("/audits", kubernetesAuditPostHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func loadAndValidateParameters() {
	tenantID = os.Getenv(tenantIDVariableName)
	clientID = os.Getenv(clientIDVariableName)
	clientSecret = os.Getenv(clientSecretVariableName)

	validateParameter(tenantID, tenantIDVariableName)
	validateParameter(clientID, clientIDVariableName)
	validateParameter(clientSecret, clientSecretVariableName)
}

func validateParameter(parameterValue string, parameterName string) {
	if parameterValue == "" {
		panic(fmt.Errorf("%s cannot be null", parameterName))
	}
}

func getCurrentDateTimeString() string {
	now := time.Now().UTC()
	formattedNow := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	return formattedNow
}
