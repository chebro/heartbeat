# systat

a simple webpage to monitor system status that is based on a client-server model where,

- client(s) send periodic pings with basic status info
- server updates the client info after each ping and renders webpage

# run on localhost

build and start server
```bash
go build -ldflags "-X main.gitCommitHash=$(git rev-parse --short HEAD)" -o systat
./systat
```

build and start client
```bash
cd client
go build -o client
./client -h localhost -p 8080
```
