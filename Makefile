NAME=zendesk_dataset
BINPATH=bin

build:
	GOOS=windows GOARCH=386   go build -o $(BINPATH)/$(NAME)_windows_x86.exe main.go
	GOOS=windows GOARCH=amd64 go build -o $(BINPATH)/$(NAME)_windows_x64.exe main.go
	GOOS=darwin  GOARCH=amd64 go build -o $(BINPATH)/$(NAME)_osx_x64 main.go
	GOOS=linux   GOARCH=386   go build -o $(BINPATH)/$(NAME)_linux_x86 main.go
	GOOS=linux   GOARCH=amd64 go build -o $(BINPATH)/$(NAME)_linux_x64 main.go

test:
	go test -v ./...
