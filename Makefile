build:
	go build -o easybackup cmd/*

install: build
	mv ./easybackup $$GOPATH/bin

clean:
	rm -f easybackup
