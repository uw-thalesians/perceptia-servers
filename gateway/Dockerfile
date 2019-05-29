FROM golang:1.12 as builder
WORKDIR /perceptia-servers/gateway/gateway/
COPY ./gateway .
RUN env GOOS=linux GOARCH=386 go build -o gateway .

FROM alpine:3.9
RUN apk add --no-cache ca-certificates
COPY --from=builder /perceptia-servers/gateway/gateway/gateway /gateway
EXPOSE 443
ENTRYPOINT ["/gateway"]