package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	graphapi "github.com/carlsoncoder/audit-webhook-go/graphapi"
	kubernetestypes "github.com/carlsoncoder/audit-webhook-go/kubernetestypes"
)

var (
	graphAPIClient      *graphapi.GraphAPIClient
	userTenantURLPrefix string
)

func kubernetesaudits(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var eventList kubernetestypes.EventList

	err = json.Unmarshal(body, &eventList)
	if err != nil {
		panic(err)
	}

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
	// load the parameters from OS environment variables and initialize our graphAPIClient class
	tenantID := os.Getenv("TENANT_ID")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	graphAPIClient = graphapi.NewGraphAPIClient(tenantID, clientID, clientSecret)
	userTenantURLPrefix = fmt.Sprintf("https://sts.windows.net/%s/#", tenantID)

	http.HandleFunc("/audits", kubernetesaudits)
	log.Fatal(http.ListenAndServe(":80", nil))
}
