FROM golang:1.11.5 as builder
WORKDIR /perceptia-servers/gateway/gateway/
COPY ./gateway .
RUN env GOOS=linux GOARCH=386 go build -o gateway .

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /perceptia-servers/gateway/gateway/gateway /gateway
EXPOSE 443
ENTRYPOINT ["/gateway"]