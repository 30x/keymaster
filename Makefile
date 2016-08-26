#Simply build for linux.  This needs to have testing and validation included

release: update-deps compile-linux

update-deps:
	glide install

compile-linux:
	GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o build/keymaster .
