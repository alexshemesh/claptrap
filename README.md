# Get the code

go get github.com/alexshemesh/claptrap

# Go to folder

cd $GOPATH/src/github.com/alexshemesh/claptrap

# Dependencies GLIDE

We use glide https://github.com/Masterminds/glide to manage dependencies
after that:

glide up

# Test
go test --cover $(go list ./... | grep -v '/vendor/')

# Build

mkdir build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=1.0.0" -o build/claptrap-linux-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.Version=1.0.0" -o build/claptrap-darwin-amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.Version=1.0.0" -o build/claptrap-windows-amd64
 


