FROM hypriot/rpi-alpine
MAINTAINER Evan Lezar <evanlezar@gmail.com>

RUN apk add --no-cache ca-certificates

COPY go-faster.rpi /usr/local/bin
ENTRYPOINT ["/usr/local/bin/go-faster.rpi"]