# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go build -o audit-webhook-go

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/audit-webhook-go /app/
ENTRYPOINT ./audit-webhook-go