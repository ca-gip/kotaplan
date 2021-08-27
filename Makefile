.PHONY: clean

REPO= github.com/ca-gip/dploy
NAME= kotaplan

dependency:
	go mod download

linux:
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -v -o ./build/linux_amd64 -i $(GOPATH)/src/$(REPO)/main.go

darwin:
	GOOS=darwin CGO_ENABLED=0 GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -v -o ./build/darwin_amd64 -i $(GOPATH)/src/$(REPO)/main.go

build: linux darwin