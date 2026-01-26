FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o BingPaper .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/BingPaper .
RUN mkdir -p data
COPY --from=builder /app/config.example.yaml ./data/config.yaml
COPY --from=builder /app/web ./web

EXPOSE 8080
ENTRYPOINT ["./BingPaper"]
