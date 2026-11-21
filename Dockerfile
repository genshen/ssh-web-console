# build method: just run `docker build --rm -t genshen/ssh-web-console .`

# build frontend code
# Node 16: the frontend (web submodule) is react-scripts 4, whose PostCSS does not
# load on Node >= 17. Only the static build output leaves this stage, so no Node
# runtime ships in the image. Moving the frontend off react-scripts 4 is the real fix.
FROM node:16-alpine AS frontend-builder

COPY web web-console/

RUN cd web-console \
    && yarn install \
    && yarn build


FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
# the frontend is embedded into the binary from ./build (see main.go)
COPY --from=frontend-builder web-console/build ./build/

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/ssh-web-console .

## copy binary
FROM alpine:3.24

ARG HOME="/home/web"

# pick up Alpine security fixes published since the base image was built
RUN apk upgrade --no-cache \
    && adduser -D web -h ${HOME}

COPY --from=builder --chown=web /out/ssh-web-console ${HOME}/ssh-web-console

WORKDIR ${HOME}
USER web

VOLUME ["${HOME}/conf"]

CMD ["./ssh-web-console"]
