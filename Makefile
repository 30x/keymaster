#Simply build for linux.  This needs to have testing and validation included

#This is because we have a test directory, and make thinks it doesn't have to do anything
.PHONY: test

release: update-deps test compile-linux

test:
	go test $$(glide novendor)

update-deps:
	glide install

compile-linux:
	GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o build/keymaster .
