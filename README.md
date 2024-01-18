# DecTimer
A silly Discord bot for tracking when Declan is back off ship. Written in Go.

# Features to add
- Live tracking ship via https://aisstream.io/

# Dev commands
 go run .\DecTimer.go -t

# Docker commands
docker build --rm -t dectimer:alpha .

docker run -d -p 8080:8081 --name decobot dectimer:alpha
