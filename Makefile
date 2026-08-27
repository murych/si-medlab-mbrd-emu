SHELL := /bin/sh

GO ?= go
BUILD_DIR := build
THIRD_PARTY_BUILD := $(BUILD_DIR)/libsrvint
SRVINT_SOURCE := third_party/libsrvint
SRVINT_ARTIFACTS := internal/srvint/third_party
EMU := $(BUILD_DIR)/mbrd-emu

.PHONY: 3rdparty build test clean

3rdparty:
	git submodule update --init --recursive
	cmake -S $(SRVINT_SOURCE) -B $(THIRD_PARTY_BUILD) -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=OFF
	cmake --build $(THIRD_PARTY_BUILD)
	mkdir -p $(SRVINT_ARTIFACTS)/include $(SRVINT_ARTIFACTS)/lib
	cp $(SRVINT_SOURCE)/src/srvint.h $(SRVINT_ARTIFACTS)/include/srvint.h
	cp $(THIRD_PARTY_BUILD)/export/libsrvint/libsrvint_export.h $(SRVINT_ARTIFACTS)/include/libsrvint_export.h
	cp $(THIRD_PARTY_BUILD)/libsrvint.a $(SRVINT_ARTIFACTS)/lib/libsrvint.a

build: 3rdparty
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build -o $(EMU) ./cmd/mbrd-emu

test: 3rdparty
	CGO_ENABLED=1 $(GO) test ./...

clean:
	rm -rf $(BUILD_DIR) $(SRVINT_ARTIFACTS)
