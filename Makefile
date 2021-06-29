build:
	go mod tidy
	go build

test:
	rm -rf testarea
	make install
	mkdir testarea
	go test -race -coverprofile=coverage.txt -covermode=atomic
	rm -rf testarea

install:
	@cp -v themes/*.yaml ~/.config/ffcss/themes/
	@cp -v ffcss ~/.local/bin/ffcss

format:
	gofmt -s -w **.go
