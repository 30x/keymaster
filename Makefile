#Simply build for linux.  This needs to have testing and validation included

compile-linux:
	GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o build/keymaster .
