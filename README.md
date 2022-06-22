# streambot

## Requirements
**ZeroMQ**
```shell
apt-get install libzmq3-dev
```

## Build
### Options
Prevent interrupted system calls [see docs](https://pkg.go.dev/github.com/pebbe/zmq4#section-documentation).
```
GODEBUG=asyncpreemptoff=1
```
