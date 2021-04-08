FROM golang:1.15
WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o local-php-security-checker .

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/local-php-security-checker /
ENTRYPOINT [ "/local-php-security-checker" ]