FROM golang:1.15 as proxy
ENV GO111MODULE=on
WORKDIR /go/src/proxy
COPY . /go/src/proxy
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64   go build  ./main.go

FROM alpine
WORKDIR /app
COPY --from=proxy /go/src/proxy /app
RUN chmod -R +x  .
EXPOSE 8080/tcp
ENTRYPOINT [ "/app/main", "--proto=http"]