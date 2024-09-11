ARCH ?= x86_64

BIN = build/bin/brainstack
APPIMAGETOOL = build/appimagetool-$(ARCH).AppImage
APPIMAGE = Brainstack-$(ARCH).AppImage

$(BIN):
	wails build

appimage: build/$(APPIMAGE)

$(APPIMAGETOOL):
	curl -L -o $(APPIMAGETOOL) https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-$(ARCH).AppImage
	chmod +x $(APPIMAGETOOL)

build/$(APPIMAGE): $(BIN) $(APPIMAGETOOL)
	cp $(BIN) appdir/AppRun
	./$(APPIMAGETOOL) appdir
	chmod +x $(APPIMAGE)
	mv $(APPIMAGE) build

clean:
	rm build/$(APPIMAGE)

clean-bin: clean
	rm $(BIN)

.PHONY: clean clean-bin
