# build stage
FROM golang:alpine AS build-env
ADD . /go/src/github.com/carlsoncoder/audit-webhook-go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 cd /go/src/github.com/carlsoncoder/audit-webhook-go && go build -o audit-webhook-go

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/carlsoncoder/audit-webhook-go/audit-webhook-go /app/
ENTRYPOINT ./audit-webhook-go
