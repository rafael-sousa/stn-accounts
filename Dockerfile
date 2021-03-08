FROM golang:1.15-alpine

RUN mkdir -p /go/src/github.com/rafael-sousa/stn-accounts
WORKDIR /go/src/github.com/rafael-sousa/stn-accounts
ADD . .
RUN mkdir /usr/local/app
RUN go build -o /usr/local/app/main /go/src/github.com/rafael-sousa/stn-accounts/cmd/rip-go/main.go
RUN rm -rf /go/src/github.com/rafael-sousa/stn-accounts
CMD ["/usr/local/app/main"]