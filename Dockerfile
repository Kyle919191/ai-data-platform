FROM golang:1.24

WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd

RUN go build -o ingestion-gateway ./cmd/ingestion-gateway

EXPOSE 8080

CMD ["./ingestion-gateway"]