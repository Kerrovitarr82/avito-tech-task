FROM golang:1.25.4

WORKDIR /app
COPY ./ ./
RUN go mod download && go mod verify
RUN go build -o pr-reviewer ./cmd/app/main.go
CMD ["./pr-reviewer"]