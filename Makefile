PREFIX ?= /opt/brainstack
BINPREFIX ?= $(PREFIX)/bin

DESKTOPPREFIX ?= $(HOME)/.local/share/applications

BIN = build/bin/brainstack

all: install

$(BIN):
	wails build

install: $(BIN)
	sudo mkdir -p $(BINPREFIX)
	sudo cp $(BIN) $(BINPREFIX)/brainstack

	mkdir -p $(DESKTOPPREFIX)
	cp brainstack.desktop $(DESKTOPPREFIX)/brainstack.desktop

.PHONY: all install
