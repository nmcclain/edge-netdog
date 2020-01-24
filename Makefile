
all: linux pi pizero osx mipsle mips win32 win64

linux:
	GOOS=linux GOARCH=amd64 go build -o build/edge-netdog_linux-amd64

pi:
	GOOS=linux GOARCH=arm GOARM=7 go build -o build/edge-netdog_linux-arm7

pizero:
	GOOS=linux GOARCH=arm GOARM=6 go build -o build/edge-netdog_linux-arm6

osx:
	GOOS=darwin GOARCH=amd64 go build -o build/edge-netdog_darwin-amd64

mipsle:
	GOOS=linux GOARCH=mipsle go build -o build/edge-netdog_linux-mipsle

mips:
	GOOS=linux GOARCH=mips go build -o build/edge-netdog_linux-mips

win32:
	GOOS=windows GOARCH=386 go build -o build/edge-netdog_windows-386

win64:
	GOOS=windows GOARCH=amd64 go build -o build/edge-netdog_windows-amd64
