# build method: just run `docker build --rm --build-arg -t genshen/ssh-web-console .`

# build frontend code
FROM node:10-alpine AS frontend-builder

ARG FRONTEND_REP="https://github.com/genshen/webConsole.git"
ARG FRONTEND_VERSION="master"

RUN apk add --no-cache git \
    && cd /  \
    && git clone ${FRONTEND_REP} web-console \
    && cd web-console \
    && git checkout ${FRONTEND_VERSION} \
    && yarn install \
    && yarn build


FROM golang:1.12.3-alpine AS builder

# set to 'on' if using go module
ENV GO111MODULE=on
ARG STATIC_DIR=dist

RUN apk add --no-cache git \
    && go get -u github.com/rakyll/statik

COPY ./  /go/src/github.com/genshen/ssh-web-console/
COPY --from=frontend-builder /web-console/dist /go/src/github.com/genshen/ssh-web-console/${STATIC_DIR}/

RUN cd ./src/github.com/genshen/ssh-web-console/ \
    && statik -src=${STATIC_DIR} \
    && go build \
    && go install

## copy binary
FROM alpine:latest
ARG USRR="web"
ARG HOME="/home/web"

RUN adduser -D ${USRR}

COPY --from=builder --chown=web /go/bin/ssh-web-console ${HOME}/ssh-web-console

WORKDIR ${HOME}
USER ${USER}

VOLUME ["${HOME}/conf", "${HOME}/views"]

# fixme still using root user.
CMD ["./ssh-web-console"]
