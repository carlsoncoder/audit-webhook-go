# audit-webhook-go

An implementation of an audit webhook to be used by the kubernetes api-server, written in Go

TODO ITEMS:

* "graphapiclient.go" -  Shouldn't be using "panic"...instead log the error and continue on (return null??) so we don't crash the program
* When iterating through "eventList.Events", we shouldn't call the graphapiclient multiple times for the same user
* Find some way to cache the access_token in graphapiclient so we aren't repeatedly calling it when we don't have to
* TEST IT ALL OUT!
* FINISH UPDATING THIS README - include example with curl & test JSON file (such as "curl -X POST -H 'Content-Type: application/json' -d @sampleKubernetesAuditPostData.json http://localhost/audits").  include notes about setting environment variables and changing your values!