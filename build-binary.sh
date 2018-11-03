#bin/bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/statefulset-operator cmd/statefulset-operator/main.go