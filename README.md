# sshWebConsole
you can connect to your linux machine by ssh in your browser.

## Dependency
```bash
go get -u github.com/kardianos/govendor  # ues [govendor](https://github.com/kardianos/govendor) to manager dependency
```

## How to build
1. clone the repository [webConsole](https://github.com/genshen/webConsole) to any directory (example:/home/foo/webConsole) you like,and follow its README to build the frontend code.
2. copy the built files to present project,and edit configure file:
   ```bash
   cp /home/foo/webConsole/dist/static/  ./static/
   cp /home/foo/webConsole/dist/index.html  ./views/index.html
   cp conf/config.yaml.example conf/config.yaml
   vi conf/config.yaml  # edit configure file
   ```
3. get Dependency(run ***govendor sync***) and then run:***go build main.go*** to build present project.
4. run: ./main ,and than you can enjoy it in your browser.
![](./Screenshots/shot1.png)

## Screenshots
![](./Screenshots/shot2.png)
![](./Screenshots/shot3.png)
![](./Screenshots/shot4.png)

# Related Works
https://github.com/shibingli/webconsole
