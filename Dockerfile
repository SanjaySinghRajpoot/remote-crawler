FROM golang:1.19-alpine3.17 AS builder
# WORKDIR /app
# COPY . .
# RUN go build -o main main.go

# FROM alpine:3.13
# WORKDIR /app
# COPY --from=builder /app/main .

# EXPOSE 8080
# CMD ["/app/main"]

# set working directory
WORKDIR /app

# Copy the source code
COPY . . 

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o api .

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./api"]