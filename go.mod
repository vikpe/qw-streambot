module github.com/vikpe/streambot

go 1.18

require (
	github.com/fatih/color v1.13.0
	github.com/go-resty/resty/v2 v2.7.0
	github.com/goccy/go-json v0.9.7
	github.com/joho/godotenv v1.4.0
	github.com/pebbe/zmq4 v1.2.9
	github.com/stretchr/testify v1.7.4
	github.com/vikpe/serverstat v0.1.81
	golang.org/x/exp v0.0.0-20220613132600-b0d781184e0d
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/ssoroka/slice v0.0.0-20220402005549-78f0cea3df8b // indirect
	github.com/vikpe/udpclient v0.1.3 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sys v0.0.0-20220627191245-f75cf1eec38b // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/vikpe/serverstat => ../serverstat
