# Use official Golang image to build the application
FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main .

# Use an Alpine image for the final image to reduce size
FROM alpine:latest
RUN apk add --no-cache bash

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
