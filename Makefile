sources := $(wildcard *.go)

gobw: $(sources)
	go build

.PHONY: install
install: gobw
	install ./gobw $${HOME}/.local/bin
