test:
	go test -v ./... -cover
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o sc-cli cmd/main.go
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o sc-cli cmd/main.go
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o sc-cli.exe cmd/main.go
build-linux-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -o sc-cli cmd/main.go
