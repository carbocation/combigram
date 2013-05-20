GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o combigram.linux  *.go
GOOS=darwin GOARCH=amd64 go build -o combigram.osx  *.go

