build:
	@cp -v themes/*.yaml ~/.config/ffcss/themes/
	go mod tidy
	go build

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic

install:
	@cp -v ffcss ~/.local/bin/ffcss

format:
	gofmt -s -w **.go
