# build stage
FROM golang:alpine AS build
ADD . /go/src/github.com/carlsoncoder/audit-webhook-go
RUN apk update && apk add git
# NOTE: "go get all" seems to fail and return a non-zero exit code, which would make docker build fail
# However, all the packages WE need are still downloaded.  By OR'in the go get all command with /bin/true, we ensure
# the command will NEVER fail...this is good in our case, but if it failed for any other reason (say, network issues), this would cause problems
# However, if we really didn't get what we need - the go build command would likely fail anyway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 cd /go/src/github.com/carlsoncoder/audit-webhook-go && go get all || /bin/true && go build -o audit-webhook-go

# final stage
FROM alpine
WORKDIR /app
RUN apk update && apk add bash curl
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/carlsoncoder/audit-webhook-go/audit-webhook-go /app/
ENTRYPOINT ./audit-webhook-go