# go-webserver

A simple Go repository to show how to reference external packages and serve web requests in Go.

Once running, you can call it using curl from the root directory like so to see a sample POST:

    curl -X POST -H 'Content-Type: application/json' -d @sampleUserPostData.json http://localhost/usertest

Or, for another example:    

    curl -X POST -H 'Content-Type: application/json' -d @sampleKubernetesAuditPostData.json http://localhost/kubernetestest