env CGO_ENABLED=0 
env GOOS=linux
env GOARCH=amd64
go build -tags auto