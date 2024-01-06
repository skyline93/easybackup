build:
	go build -o easybackup cmd/*

install:
	mv easybackup ~/go/bin

clean:
	rm -f easybackup
