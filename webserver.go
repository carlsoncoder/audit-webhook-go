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
