# ssh-web-console
you can connect to your linux machine by ssh in your browser.

## Quick start
```bash
$ docker pull genshen/ssh-web-console:latest
# docker build --build-arg GOMODULE=on -t genshen/ssh-web-console . # or build docker image on your own machine
$ docker run -v ${PWD}/conf:/home/web/conf -p 2222:2222 --rm genshen/ssh-web-console
```

Open your browser, visit `http://localhost:2222`. Enjoy it!

**note**: To run docker container, make sure config.yaml file is in directory ${PWD}/conf
## Build & Run
make sure you go version is not less than 1.11

### build frontend
```bash
$ cd /tmp
$ git clone https://github.com/genshen/webConsole web-console
$ cd web-console
$ yarn install
$ yarn build
```

### build go
```bash
$ go get github.com/rakyll/statik
$ cp -r /tmp/web-console/dist  ./dist
$ statik dist  # use statik tool to convert files in 'dist' dir to go code, and compile into binary.
$ export GO111MODULE=on # for go 1.11.x
$ go build
```

## Run
run: `./ssh-web-console` (you should check directory specified by `soft_static_dir` in conf/config.yaml exists), 
and than you can enjoy it in your browser by visiting `http://localhost:2222`.

![](./Screenshots/shot1.png)

## Screenshots
![](./Screenshots/shot2.png)
![](./Screenshots/shot3.png)
![](./Screenshots/shot4.png)

# Related Works
https://github.com/shibingli/webconsole
