build:
	cd ./src && go build -o ../dist/check_ping_go *.go

install:
	cp -v ./dist/check_ping_go /usr/local/nagios/libexec/
