
release: check-ver clean
	goxc -wlc -pv=$(ver)
	goxc -wlc default publish-github -apikey=$(ghtoken)
	goxc -d bin -bc "linux"

clean:
	rm -Rf bin

test: 
	go test -v ./...

dep:
	dep ensure
		
get-tools: 
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/laher/goxc

check-ver:
	@ if [ "$(ver)" = "" ]; then \
		echo "'ver' is not set"; \
		exit 1; \
	fi

.PHONY: release clean test dep tools check-ver