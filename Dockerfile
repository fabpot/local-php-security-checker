FROM golang:1.13.8
WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o local-php-security-checker .

FROM alpine:3.11
COPY --from=0 /go/src/local-php-security-checker /
ENTRYPOINT [ "/local-php-security-checker" ]