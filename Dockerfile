# We can make use of Multi stage docker builds that will only copy a the requied code in the production image 
# It will make use of the existing base image

# alpine is used to get the lite version of Golang
FROM golang:alpine

RUN mkdir /app

WORKDIR /app

ADD go.mod .
ADD go.sum .


RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

ADD . .

EXPOSE 8000

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
