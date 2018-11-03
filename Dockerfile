FROM alpine:3.6

COPY ./bin/statefulset-operator /usr/local/bin/statefulset-operatorr
RUN apk update && apk add ca-certificates
RUN adduser -D statefulset-operatorr

USER statefulset-operatorr

ENTRYPOINT ["/usr/local/bin/statefulset-operatorr"]