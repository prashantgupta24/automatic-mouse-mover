COVER_PROFILE=cover.out
COVER_HTML=cover.html

all: open

build: clean
	mkdir -p -v ./bin/amm.app/Contents/Resources
	mkdir -p -v ./bin/amm.app/Contents/MacOS
	cp ./appInfo/*.plist ./bin/amm.app/Contents/Info.plist
	cp ./appInfo/*.icns ./bin/amm.app/Contents/Resources/icon.icns
	go build -o ./bin/amm.app/Contents/MacOS/amm cmd/main.go

open: build
	open ./bin

clean:
	rm -rf ./bin

start:
	go run cmd/main.go

coverage: $(COVER_HTML)

$(COVER_HTML): $(COVER_PROFILE)
	go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)

$(COVER_PROFILE):
	go test -v -failfast -race -coverprofile=$(COVER_PROFILE) ./...

vet:
	go vet $(shell glide nv)

lint:
	go list ./... | grep -v vendor | xargs -L1 golint -set_exit_status

.PHONY: build 
.PHONY: clean