FROM golang:1.17.6-alpine

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN go build -o elasticsearch-api -v ./...

CMD [ "./elasticsearch-api" ]
