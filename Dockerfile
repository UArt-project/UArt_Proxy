FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY ../go.mod go.sum ./

RUN go mod download

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o uart-proxy cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/uart-proxy .
COPY --from=builder /app/config.yaml .

EXPOSE 8000

CMD ["./uart-proxy"]
