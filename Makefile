GOMOD=win_hotkeys_daemon
build: release
debug: compdbg
release: comprel
comprel:
	go build -ldflags "-s -w -H=windowsgui" .
compdbg:
	go build -race -gcflags=all=-d=checkptr=0 .
clean:
	del $(GOMOD)*
