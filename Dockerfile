FROM alpine
MAINTAINER Evan Lezar <evanlezar@gmail.com>

RUN apk add --update \ 
    python \
    python-dev \
    py-pip \
    build-base \
    git \ 
  && pip install virtualenv \
  && pip install git+https://github.com/sivel/speedtest-cli.git \
  && apk del git \
  && rm -rf /var/cache/apk/*

ENTRYPOINT ["speedtest-cli"]