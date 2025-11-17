FROM golang:1.21-alpine

WORKDIR /app

COPY main.go .

# Initialize module and download dependencies
RUN go mod init api-service && go mod tidy

# Build the application
RUN go build -o main .

EXPOSE 8000

CMD ["./main"]