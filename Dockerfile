FROM golang:alpine3.16 as builder

WORKDIR /builder

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /goapp ./cmd/

## Deploy
FROM alpine:3.16

WORKDIR /

COPY --from=builder /goapp /goapp

EXPOSE 8080

ENTRYPOINT ["/goapp"]