FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/docs ./docs

RUN apk add --no-cache file
RUN ls -l /app && file /app/app

EXPOSE 8080

ENV PORT=8080

CMD ["./app"]