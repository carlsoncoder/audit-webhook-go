# build stage
FROM golang:alpine AS build
ADD . /go/src/github.com/carlsoncoder/audit-webhook-go
RUN apk update && apk add git
# Note: we *should* be able to do 'go get all' instead of 'go get github.com/dtzar/oms-go/oms_data_collector', but something was goofy with the PATH in the build container...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 cd /go/src/github.com/carlsoncoder/audit-webhook-go && go get github.com/dtzar/oms-go/oms_data_collector && go build -o audit-webhook-go

# final stage
FROM alpine
WORKDIR /app
RUN apk update && apk add bash curl
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/carlsoncoder/audit-webhook-go/audit-webhook-go /app/
ENTRYPOINT ./audit-webhook-go
