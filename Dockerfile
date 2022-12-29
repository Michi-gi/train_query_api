FROM golang

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

ENV GOPATH=

RUN go build

ENTRYPOINT ["./trainquery"]
