SHELL:=/bin/bash

.ONESHELL:
.PHONY: coverage

build:
	go mod tidy
	go build

tests:
# setup mocks
	$(MAKE) mocks-setup
# defer tearing down mocks (so that it runs even if the tests fail)
	trap "$(MAKE) mocks-teardown" EXIT
# run tests, with $HOME overriden to a mocked directory, and gopath re-set, otherwise the compiler freaks out.
	GOPATH=$$(go env GOPATH) HOME=testarea/home go test -race -coverprofile=coverage.txt -covermode=atomic -v
# compute code coverage
	$(MAKE) coverage

coverage:
# get the binary to convert from go coverage file to a standard one accepted by codecov & co
	go get -u github.com/jandelgado/gcov2lcov
# run it
	gcov2lcov -infile=coverage.txt -outfile=coverage/lcov.info
# remove it from the go.mod file
	go mod tidy

install:
# make necessary directories
	mkdir -p ~/.config/ffcss/themes ~/.local/bin
# copy builtin themes
	cp themes/*.yaml ~/.config/ffcss/themes/
# copy binary to some standard place that's in $PATH most of the time
	cp ffcss ~/.local/bin/ffcss

format:
	gofmt -s -w **.go

mocks-setup:
# local webserver to mock http requests, save its PID to a file so that we can kill it during teardown
	python -m http.server 8080 --bind localhost --directory mocks/ >/dev/null 2>/dev/null & echo $$! > .mockswebserverpid
# create directories that are cleaned up after use
	mkdir -p testarea/{zip-dropoff,cache,home/{.mozilla/firefox/667ekipp.default-release,.{config,cache}/ffcss,.config/ffcss/themes}}
# copy themes into mock config directory
	cp themes/*.yaml testarea/home/.config/ffcss/themes/
# copy static mocks from mocks/ to testarea/
	cp -R mocks/home/ testarea/
# create coverage directory
	mkdir -p coverage

mocks-teardown:
# remove testing artifacts
	rm -rf testarea
# kill mocks webserver
	kill -9 $$(cat .mockswebserverpid)
# remove mocks webserver PID file
	rm .mockswebserverpid


release:
# remove artifacts from previous release
	rm -rf dist/
# build & install binary
	$(MAKE) build
	$(MAKE) install
# update changelog headings
	chachacha release $$(read -p bump=; echo $$REPLY)
# remove stupid bullet points in front of <hX> tags
	sd '^- #' '#' CHANGELOG.md
# extract release notes for the new version only,
# and bump the version in go code
	./make_release_notes.rb $$(read -p bump=; echo $$REPLY)
# recompile so that the binary shows the new version when doing ffcss version
	$(MAKE) build
	$(MAKE) install
# make tagged Release commit
	git add CHANGELOG.md ffcss.go
	git commit -m "🔖 Release $$(ffcss version)"
	git tag -am v$$(ffcss version) v$$(ffcss version)
# github push, tags push, github release w/ binaries in .tar.gz, milestone close, etc.
	GITHUB_TOKEN=$$(cat .github_token) goreleaser release --release-notes release_notes.md
# remove extracted release notes
	rm release_notes.md
