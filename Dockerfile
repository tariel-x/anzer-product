FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go install -v ./...

CMD ["app"]