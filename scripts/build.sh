# Prevent interrupted system calls
# https://pkg.go.dev/github.com/pebbe/zmq4#section-documentation
GODEBUG=asyncpreemptoff=1

(cd ./cmd/proxy && go build)
(cd ./cmd/quake_manager && go build)
(cd ./cmd/twitchbot && go build)
(cd ./cmd/twitch_manager && go build)
