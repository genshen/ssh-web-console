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
$ git clone https://github.com/genshen/web-console
$ cd web-console
$ yarn install
$ yarn build
```

### build go
```bash
$ export GO111MODULE=on # for go 1.11.x
$ go build
$ cp /tmp/web-console/dist/static/  ./static/
$ cp /tmp/web-console/dist/index.html  ./views/index.html
```

## Run
run: `./ssh-web-console` ,and than you can enjoy it in your browser by visiting `http://localhost:2222`.
![](./Screenshots/shot1.png)

## Screenshots
![](./Screenshots/shot2.png)
![](./Screenshots/shot3.png)
![](./Screenshots/shot4.png)

# Related Works
https://github.com/shibingli/webconsole
