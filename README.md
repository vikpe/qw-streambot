# streambot [![Test](https://github.com/vikpe/streambot/actions/workflows/test.yml/badge.svg)](https://github.com/vikpe/streambot/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/vikpe/streambot/branch/main/graph/badge.svg)](https://codecov.io/gh/vikpe/streambot)
> Automated QuakeWorld streaming on Twitch.

## Overview
![image](https://user-images.githubusercontent.com/1616817/178285267-eade607d-8660-4b4d-9522-ab3772dde229.png)

## Requirements
**ZeroMQ**
```shell
apt-get install libzmq3-dev
```

**Access tokens**
* [Chat acess token for bot](https://twitchapps.com/tmi/)
* [General access token](https://twitchtokengenerator.com/)

## Build

```shell
./scripts/build.sh
```

### Options
(zeromq) Prevent interrupted system calls [see docs](https://pkg.go.dev/github.com/pebbe/zmq4#section-documentation).
```
GODEBUG=asyncpreemptoff=1
```
