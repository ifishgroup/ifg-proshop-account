FROM golang:1.8.3 AS build

ADD . /go/src/github.com/ifishgroup/ifg-proshop-account
WORKDIR /go/src/github.com/ifishgroup/ifg-proshop-account
RUN go get -d -v -t
RUN go test --cover -v ./...
RUN go build -v -o account-service


FROM alpine:3.6

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
CMD ["account-service"]
COPY --from=build /go/src/github.com/ifishgroup/ifg-proshop-account/account-service /usr/local/bin/account-service
RUN chmod +x /usr/local/bin/account-service
