# Start with the official Golang image as the base
FROM golang:1.22-alpine

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go get ./...
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
