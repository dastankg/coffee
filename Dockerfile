FROM golang:1.24 AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN swag init --pd -g cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./migrations/auto.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

# Use minimal image for final container
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app .


EXPOSE 8081

CMD ["sh", "-c", "./migrate && ./main"]