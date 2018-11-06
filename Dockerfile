FROM golang:latest
RUN go get -u github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/rbo13/write-it

ADD . /go/src/github.com/rbo13/write-it
COPY ./Gopkg.toml /go/src/github.com/rbo13/write-it
WORKDIR /go/src/github.com/rbo13/write-it

RUN dep ensure
RUN go test -v 
RUN go build -o main .

EXPOSE 1333
CMD ["./main"]