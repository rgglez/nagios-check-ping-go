build:
	cd ./src && go build -o ../dist/check_icmp_go *.go

install:
	cp -v ./dist/check_icmp_go /usr/local/nagios/libexec/