package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	auditextensions "github.com/carlsoncoder/audit-webhook-go/auditextensions"
	graphapi "github.com/carlsoncoder/audit-webhook-go/graphapi"
	omstypes "github.com/carlsoncoder/audit-webhook-go/omstypes"
	omsgo "github.com/dtzar/oms-go/oms_data_collector"
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
)

var (
	tenantID            string
	clientID            string
	clientSecret        string
	omsCustomerID       string
	omsSharedKey        string
	graphAPIClient      graphapi.Client
	omsLogClient        omsgo.OmsLogClient
	userTenantURLPrefix string
)

const (
	tenantIDVariableName      = "TENANT_ID"
	clientIDVariableName      = "CLIENT_ID"
	clientSecretVariableName  = "CLIENT_SECRET"
	omsCustomerIDVariableName = "OMS_CUSTOMER_ID"
	omsSharedKeyVariableName  = "OMS_SHARED_KEY"
	omsPostTimeout            = time.Second * 5
	omsLogType                = "kubernetesaudits"
)

func kubernetesAuditPostHandler(rw http.ResponseWriter, req *http.Request) {
	now := getCurrentDateTimeString()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(fmt.Sprintf("[%s] ERROR: Unable to read POST body", now))
		log.Println(fmt.Sprintf("%v", err))
		return
	}

	var eventList auditinternal.EventList

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

	for _, item := range eventList.Items {
		// We only want to log records that came from an AAD user
		if auditextensions.IsUserValidAADUser(&item.User, userTenantURLPrefix) {
			// determine user and user group properties
			userDisplayName, userPrincipalName := auditextensions.GetUserDetails(&item.User, userTenantURLPrefix, graphAPIClient)
			userGroups := auditextensions.GetUserGroups(&item.User, graphAPIClient)

			// get the source IP address(es) from the event
			sourceIPAddress := auditextensions.GetSourceIPAddress(&item)

			authorizationDecision := item.Annotations["authorization.k8s.io/decision"]
			authorizationReason := item.Annotations["authorization.k8s.io/reason"]

			// build the message to be sent to OMS
			omsMessage := &omstypes.LogMessage{
				Timestamp:             item.StageTimestamp.Time,
				RequestURI:            item.RequestURI,
				Level:                 fmt.Sprintf("%s", item.Level),
				Stage:                 fmt.Sprintf("%s", item.Stage),
				Verb:                  item.Verb,
				UserAgent:             item.UserAgent,
				SourceIPAddress:       sourceIPAddress,
				FullUserName:          item.User.Username,
				UserDisplayName:       userDisplayName,
				UserPrincipalName:     userPrincipalName,
				UserGroups:            userGroups,
				ResourceType:          item.ObjectRef.Resource,
				ResourceName:          item.ObjectRef.Name,
				ResourceNamespace:     item.ObjectRef.Namespace,
				ResponseStatus:        item.ResponseStatus.Status,
				ResponseReason:        fmt.Sprintf("%s", item.ResponseStatus.Reason),
				ResponseCode:          item.ResponseStatus.Code,
				AuthorizationDecision: authorizationDecision,
				AuthorizationReason:   authorizationReason}

			// print the message
			omsMessage.Print()

			// POST the message to OMS
			err := omsMessage.PostToOMS(omsLogClient, omsLogType)
			if err != nil {
				log.Println(fmt.Sprintf("[%s] ERROR: Unable to POST message to OMS", now))
				log.Println(fmt.Sprintf("%v", err))

				// still want to try and process the rest of the event messages!
				continue
			}
		}
	}
}

func main() {
	// load the necessary parameters from OS environment variables and validate they are present
	loadAndValidateParameters()

	// initialize our graphAPIClient with the parameters
	graphAPIClient = graphapi.NewClient(tenantID, clientID, clientSecret)
	userTenantURLPrefix = fmt.Sprintf("https://sts.windows.net/%s/#", tenantID)

	// initialize the OMS GO client with the parameters
	omsLogClient = omsgo.NewOmsLogClient(omsCustomerID, omsSharedKey, omsPostTimeout)

	// setup the handler for POSTing to the /audits endpoint
	http.HandleFunc("/audits", kubernetesAuditPostHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func loadAndValidateParameters() {
	tenantID = os.Getenv(tenantIDVariableName)
	clientID = os.Getenv(clientIDVariableName)
	clientSecret = os.Getenv(clientSecretVariableName)
	omsCustomerID = os.Getenv(omsCustomerIDVariableName)
	omsSharedKey = os.Getenv(omsSharedKeyVariableName)

	validateParameter(tenantID, tenantIDVariableName)
	validateParameter(clientID, clientIDVariableName)
	validateParameter(clientSecret, clientSecretVariableName)
	validateParameter(omsCustomerID, omsCustomerIDVariableName)
	validateParameter(omsSharedKey, omsSharedKeyVariableName)
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
