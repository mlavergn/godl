###############################################
#
# Makefile
#
###############################################

.DEFAULT_GOAL := build

VERSION := 1.1.0

ver:
	@sed -i '' 's/^const Version = "[0-9]\{1,3\}.[0-9]\{1,3\}.[0-9]\{1,3\}"/const Version = "${VERSION}"/' main.go

build:
	go build  -o godl main.go

rt:
	GOARCH=arm GOARM=5 GOOS=linux go build --ldflags "-s -w" -o godl main.go
