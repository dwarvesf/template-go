FROM golang:1.17
RUN mkdir /build
WORKDIR /build
COPY . .

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
RUN go install -v ./...
RUN go install -v github.com/rubenv/sql-migrate/sql-migrate@latest

FROM alpine:3.15.0
ARG DEFAULT_PORT
RUN apk --no-cache add ca-certificates
WORKDIR /

COPY --from=0 /go/bin/* /usr/bin/
COPY template /srv/neko/template
COPY migrations /migrations
COPY dbconfig.yml /
## config for timezone
COPY --from=0 /usr/share/zoneinfo /usr/share/zoneinfo
EXPOSE ${DEFAULT_PORT}

CMD [ "server" ]
