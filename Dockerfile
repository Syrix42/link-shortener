FROM golang:1.25.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/main.go -o ./swagger --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server ./server

EXPOSE 3000

CMD ["./server"]
