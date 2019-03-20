package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	graphapi "github.com/carlsoncoder/audit-webhook-go/graphapi"
	kubernetestypes "github.com/carlsoncoder/audit-webhook-go/kubernetestypes"
)

var (
	userFinder          *graphapi.UserFinder
	userTenantURLPrefix string
)

func kubernetesaudits(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	log.Println(string(body))
	var eventList kubernetestypes.EventList

	err = json.Unmarshal(body, &eventList)
	if err != nil {
		panic(err)
	}

	for i, event := range eventList.Events {
		// We only want to log records that came from an AAD user
		if strings.HasPrefix(event.User.UserName, userTenantURLPrefix) {
			log.Println(fmt.Sprintf("Event #%d", i+1))
			log.Println(event.StageTimestamp)
			log.Println(event.Level)
			log.Println(event.Stage)
			log.Println(event.RequestURI)
			log.Println(event.Verb)

			// determine and output user properties
			userObjectID := event.User.UserName[len(userTenantURLPrefix):]
			userDisplayName, userPrincipalName := userFinder.LoadUserValues(userObjectID)
			log.Println(userObjectID)
			log.Println(userDisplayName)
			log.Println(userPrincipalName)

			log.Println(event.GetSourceIPAddress())
			log.Println(event.ObjectRef.Name)
			log.Println(event.ObjectRef.Resource)
			log.Println(event.Annotations.Decision)
			log.Println(event.Annotations.Reason)
		}
	}
}

func main() {
	// load the parameters from OS environment variables and initialize our userFinder class
	tenantID := os.Getenv("TENANT_ID")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	userFinder = graphapi.NewUserFinder(tenantID, clientID, clientSecret)
	userTenantURLPrefix = fmt.Sprintf("https://sts.windows.net/%s/#", tenantID)

	http.HandleFunc("/audits", kubernetesaudits)
	log.Fatal(http.ListenAndServe(":80", nil))
}

// TODO: JUSTIN: Test a POST to the endpoint and ensure ALL values (including annotations) are correctly loaded, and that the user logging validation works
// TODO: JUSTIN: "userfinder.go" -  Shouldn't be using "panic"...instead log the error and continue on (return null??) so we don't crash the program
// TODO: JUSTIN: Include the "groups" that a user is part of - in the JSON this is the user.groups, which is a string[]
// TODO: JUSTIN: When iterating through "eventList.Events", we shouldn't call the userFinder multiple times for the same user
// TODO: JUSTIN: Find some way to cache the access_token in userfinder so we aren't repeatedly calling it when we don't have to
// TODO: JUSTIN: TEST IT ALL OUT!
