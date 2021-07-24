SHELL:=/bin/bash

build:
	go mod tidy
	go build

tests:
	make mocks-setup > /dev/null
	go test -race -coverprofile=coverage.txt -covermode=atomic -v
	go get -u github.com/jandelgado/gcov2lcov
	gcov2lcov -infile=coverage.txt -outfile=coverage/lcov.info
	go mod tidy
	make mocks-teardown > /dev/null

install:
	mkdir -p ~/.config/ffcss/themes ~/.local/bin
	@cp -v themes/*.yaml ~/.config/ffcss/themes/
	@cp -v ffcss ~/.local/bin/ffcss

format:
	gofmt -s -w **.go

mocks-setup:
# this is a code smell! ↓
	mkdir -p ~/.config/ffcss/themes ~/.local/bin
	@cp -v themes/*.yaml ~/.config/ffcss/themes/
	mkdir -p mocks/{zip-dropoff,cache-directory,homedir/.mozilla/firefox/667ekipp.default-release} testarea
	mkdir -p coverage

mocks-teardown:
	rm -rf mocks/{zip-dropoff,cache-directory,homedir} testarea

release:
	rm -rf dist/
	make > /dev/null
	make install > /dev/null
	sd -- '^  +[*-] ' '- * ' CHANGELOG.md
	chachacha release $$(read -p bump=; echo $$REPLY)
	sd '^([*-] ){2}' '  * ' CHANGELOG.md
	sd '^- #' '#' CHANGELOG.md
	./make_release_notes.rb $$(read -p bump=; echo $$REPLY)
	make > /dev/null
	make install > /dev/null
	git add CHANGELOG.md ffcss.go
	git commit -m "🔖 Release $$(ffcss version)"
	git tag -am v$$(ffcss version) v$$(ffcss version)
	GITHUB_TOKEN=$$(cat .github_token) goreleaser release --release-notes release_notes.md
	rm release_notes.md
