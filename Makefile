###############################################
#
# Makefile
#
###############################################

 build:
	go build  -o godl main.go

rt:
	GOARCH=arm GOARM=5 GOOS=linux go build --ldflags "-s -w" -o godl main.go
