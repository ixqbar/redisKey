TARGET=redisKey

all: linux mac win

linux: 
	cd src && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../bin/redisKey_${@}

mac:
	cd src && GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../bin/redisKey_${@}

win:
	cd src && GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../bin/${TARGET}.exe

clean:
	rm -rf ./bin/${TARGET}_*	
