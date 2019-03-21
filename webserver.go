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
	now := time.Now().UTC()
	formattedNow := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	log.Println(fmt.Sprintf("[%s] Processing POST request...", formattedNow))
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Error reading POST!")
		log.Println(fmt.Sprintf("%v", err))
		return
	}

	var eventList kubernetestypes.EventList

	err = json.Unmarshal(body, &eventList)
	if err != nil {
		log.Println("Unable to parse POST to JSON")
		log.Println(fmt.Sprintf("%v", err))
		log.Println("Full POST Body:")
		log.Println(string(body[:]))
		return
	}

	// error handling function to capture any errors that occur in the processing in the rest of this method
	defer func() {
		err := recover()
		if err != nil {
			log.Println("ERROR PROCESSING KUBERNETES AUDIT RECORDS:")
			log.Println(fmt.Sprintf("%v", err))
			log.Println("Full POST Body:")
			log.Println(string(body[:]))
		}
	}()

	for i, event := range eventList.Events {
		// We only want to log records that came from an AAD user
		if event.User.IsUserValidAADUser(userTenantURLPrefix) {
			log.Println(fmt.Sprintf("Event #%d", i+1))
			log.Println(fmt.Sprintf("Timestamp: %s", event.StageTimestamp))
			log.Println(fmt.Sprintf("Level: %s", event.Level))
			log.Println(fmt.Sprintf("Stage: %s", event.Stage))
			log.Println(fmt.Sprintf("Requst URI: %s", event.RequestURI))
			log.Println(fmt.Sprintf("Verb: %s", event.Verb))

			// determine and output user properties
			userObjectID, userDisplayName, userPrincipalName := event.User.GetUserDetails(userTenantURLPrefix, graphAPIClient)
			log.Println(fmt.Sprintf("User Object ID: %s", userObjectID))
			log.Println(fmt.Sprintf("User Display Name: %s", userDisplayName))
			log.Println(fmt.Sprintf("User Principal Name: %s", userPrincipalName))

			log.Println(fmt.Sprintf("Source IP: %s", event.GetSourceIPAddress()))
			log.Println(fmt.Sprintf("User Groups: %s", event.User.GetUserGroups(graphAPIClient)))
			log.Println(fmt.Sprintf("ObjectRef Name: %s", event.ObjectRef.Name))
			log.Println(fmt.Sprintf("ObjectRef Resource: %s", event.ObjectRef.Resource))
			log.Println(fmt.Sprintf("Authorization Decision: %s", event.Annotations.Decision))
			log.Println(fmt.Sprintf("Authorization Reason: %s", event.Annotations.Reason))
		}
	}
}

func main() {
	loadParameters()

	graphAPIClient = graphapi.NewClient(tenantID, clientID, clientSecret)
	userTenantURLPrefix = fmt.Sprintf("https://sts.windows.net/%s/#", tenantID)

	http.HandleFunc("/audits", kubernetesAuditPostHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func loadParameters() {
	// load the parameters from OS environment variables and initialize our graphAPIClient class
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
