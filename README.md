# audit-webhook-go

An implementation of an audit webhook to be used by the kubernetes api-server, written in Go

In order to run this, you need an Azure AD application registration setup with the appropriate Graph API permissions. Instructions on how to configure this can be found 
[here](https://docs.microsoft.com/en-us/graph/auth-v2-service?toc=./ref/toc.json&view=graph-rest-1.0).

Additionally, you need to export/define several environment variables:

* TENANT_ID
  * Your Azure AD tenant ID
* CLIENT_ID
  * The client_id of your Azure application registration
* CLIENT_SECRET
  * The client_secret to be able to access on behalf of your Azure application registration

If you don't want to set these up as environment variables, you can optionally just hard-code them in loadParameters function in webserver.go

You can test it via CURL with the following command from the root directory of this repository:

    curl -X POST -H 'Content-Type: application/json' -d @sampleKubernetesAuditPostData.json http://localhost/audits
    
However - please note that you won't actually get any Azure AD User or Group information if you use that file.  You will want to modify the "user.username" and "user.groups" values to match up with actual user and group object ID's in your Azure AD tenant if you want to see the calls to Graph API work.

There is also a multi-stage build Dockerfile in the repo that you can use to build a Docker image with this binary.   You can create it with the following command:

    docker build -t audit-webhook-go .
