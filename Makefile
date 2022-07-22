SRC := $(shell find . -type f -name '*.go')

build: $(SRC)
	go build -o bin/spider github.com/Sora233/DDC/spider
	go build -o bin/server github.com/Sora233/DDC/server

