# streambot [![Test](https://github.com/vikpe/streambot/actions/workflows/test.yml/badge.svg)](https://github.com/vikpe/streambot/actions/workflows/test.yml)

## Requirements
**ZeroMQ**
```shell
apt-get install libzmq3-dev
```

## Twitch
**Generate access tokens**
* [Chat acess token for bot](https://twitchapps.com/tmi/)
* [General access token](https://twitchtokengenerator.com/)

## Build
### Options
Prevent interrupted system calls [see docs](https://pkg.go.dev/github.com/pebbe/zmq4#section-documentation).
```
GODEBUG=asyncpreemptoff=1
```
