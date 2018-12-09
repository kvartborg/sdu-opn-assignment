FROM golang:1.11.1 as builder
WORKDIR /go/src/take-home
COPY . /go/src/take-home
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go get github.com/lib/pq
RUN go build -o server

FROM alpine
COPY --from=builder /go/src/take-home/server /usr/local/bin
COPY --from=builder /go/src/take-home/views /var/www/
ENV HTTP_VIEWS=/var/www/
EXPOSE 8080
CMD ["server"]
