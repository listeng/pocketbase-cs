cd ui && npm run build && cd .. && go build -o pocketbase-darwin-amd64 examples/base/main.go
mv pocketbase-darwin-amd64 /Users/tls/git/fastscreen/pb-builder/


# GOOS=linux go build -o pocketbase-linux-amd64 examples/base/main.go