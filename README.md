# ssh-web-console

you can connect to your linux machine by ssh in your browser.

![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/genshen/ssh-web-console?logo=docker&sort=date)
![Docker Image Version (latest semver)](https://img.shields.io/docker/v/genshen/ssh-web-console?sort=semver&logo=docker)
![Docker Pulls](https://img.shields.io/docker/pulls/genshen/ssh-web-console?logo=docker)

## Quick start

```bash
$ docker pull genshen/ssh-web-console:latest
# docker build --build-arg GOMODULE=on -t genshen/ssh-web-console . # or build docker image on your own machine
$ docker run -v ${PWD}/conf:/home/web/conf -p 2222:2222 --rm genshen/ssh-web-console
```

or using Docker Compose:

```bash
$ git clone https://github.com/genshen/ssh-web-console.git
$ cd ssh-web-console
$ docker-compose up -d
```

Open your browser, visit `http://localhost:2222`. Enjoy it!

**NOTE**: To run docker container, make sure config.yaml file is in directory ${PWD}/conf

## Build & Run

Make sure your Go version is 1.11 or later

### Clone

```bash
git clone --recurse-submodules https://github.com/genshen/ssh-web-console.git
cd ssh-web-console
```

### Build frontend

```bash
cd web
yarn install
yarn build
cd ../
```

### Build go

```bash
go get github.com/rakyll/statik
statik --src=web/build  # use statik tool to convert files in 'web/build' dir to go code, and compile into binary.
export GO111MODULE=on # for go 1.11.x
go build
```

## Run

run: `./ssh-web-console`, and than you can enjoy it in your browser by visiting `http://localhost:2222`.

## Screenshots

![](./Screenshots/shot2.png)
![](./Screenshots/shot3.png)
![](./Screenshots/shot4.png)

## Related Works

- [shibingli/webconsole](https://github.com/shibingli/webconsole)
