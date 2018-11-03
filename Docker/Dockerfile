FROM golang:1.9-alpine

#ADD . /Users/venkatesh.kristipati/go-projects
WORKDIR /Users/venkatesh.kristipati/go-projects
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/statefulset-operator cmd/statefulset-operator/main.go

FROM alpine:3.6

COPY --from=0 /tmp/hpa-operator /usr/local/bin/statefulset-operator
RUN apk update && apk add ca-certificates
RUN adduser -D hpa-operator

USER hpa-operator

ENTRYPOINT ["/usr/local/bin/statefulset-operator"]