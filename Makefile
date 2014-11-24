
build:
	gox -os="linux darwin" -arch=amd64

clean:
	git clean -fd

.PHONY: clean