FROM golang:1.12 as builder

ENV GO111MODULE=on
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./... && go build ./... && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' -o /go/bin/cli .

# ========
FROM alpine:3.9
WORKDIR /root/

COPY --from=builder /go/bin/cli .

ENTRYPOINT [ "/root/cli" ]
