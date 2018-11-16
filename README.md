# ssh-web-console
you can connect to your linux machine by ssh in your browser.

##Quick start
```bash
$ docker build --build-arg GOMODULE=on -t genshen/ssh-web-console .
$ docker run -v ${PWD}/conf:/home/web/conf --rm genshen/ssh-web-console
```

## Build & Run
make sure you go version is not less than 1.11

### build frontend
```bash
$ cd /tmp;
$ git clone https://github.com/genshen/webConsole web-console
$ cd web-console
$ yarn install
$ yarn build
```

### build go
```bash
$ go get github.com/rakyll/statik
$ cp -r /tmp/web-console/dist  ./dist
$ statik dist  # use statik tool to convert files in 'dist' dir to go code, and compile to binary.
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
