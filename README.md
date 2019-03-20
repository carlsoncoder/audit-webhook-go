# audit-webhook-go

An implementation of an audit webhook to be used by the kubernetes api-server, written in Go

TODO ITEMS:
    JUSTIN: Test a POST to the endpoint and ensure ALL values (including annotations) are correctly loaded, and that the user logging validation works
    JUSTIN: "userfinder.go" -  Shouldn't be using "panic"...instead log the error and continue on (return null??) so we don't crash the program
    JUSTIN: Include the "groups" that a user is part of - in the JSON this is the user.groups, which is a string[]
    JUSTIN: When iterating through "eventList.Events", we shouldn't call the userFinder multiple times for the same user
    JUSTIN: Find some way to cache the access_token in userfinder so we aren't repeatedly calling it when we don't have to
    JUSTIN: Make a test JSON file to test it out and add it to the repo
    JUSTIN: TEST IT ALL OUT!
    JUSTIN: FINISH UPDATING THIS README - include example with curl & test JSON file (such as "curl -X POST -H 'Content-Type: application/json' -d @sampleKubernetesAuditPostData.json http://localhost/audits")