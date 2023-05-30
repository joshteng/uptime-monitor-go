FROM golang:1.19-buster

COPY go.mod /app/go.mod
COPY go.sum /app/go.sum

WORKDIR /app

RUN go mod download

COPY . /app

RUN go build -o uptime-monitor

# EXPOSE 8080

CMD ["./uptime-monitor"]

# docker build --tag uptime-monitor:latest .
# docker run --name=uptime-monitor uptime-monitor
