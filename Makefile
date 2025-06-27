################################
BINOUT  ?= cmbxip
GOTAGS  ?= gtk_3_24
VERNUM   = 0,2,3,1
VERSTR   = 0.2.3.1_$(shell date -u +%Y.%m.%d_%H:%M)
LDFLAGS += -s -w -X nikeron/cmbxip/config.programversion=$(VERSTR)
GOFLAGS  = -v -x -tags "$(GOTAGS)" -trimpath -buildmode=exe
GOPROPS  = GOROOT_FINAL=/ CGO_ENABLED=1 GOMIPS=softfloat
################################

all: windows_amd64 linux_386 linux_amd64 linux_arm linux_arm64 linux_mipsle

generate:

test: generate
	$(GOPROPS) \
	go test -ldflags="$(LDFLAGS)" $(GOFLAGS)

current: generate
	$(GOPROPS) \
	go build -ldflags="$(LDFLAGS)" $(GOFLAGS) -o $(BINOUT)

linux_386 linux_amd64 linux_arm linux_arm64 linux_mipsle:linux_%: generate
	$(GOPROPS) \
	GOOS=linux GOARCH=$* \
	go build -ldflags="$(LDFLAGS)" $(GOFLAGS) -o $(BINOUT)_linux_$*

windows_amd64:windows_%: generate
	x86_64-w64-mingw32.static-windres \
		-DDEF_VERNUM="$(VERNUM)" -DDEF_VERSTR=\\\"$(VERSTR)\\\" \
		-DDEF_ORIGFILENAME=\\\"$(BINOUT)_$*.exe\\\" \
		-c 65001 -o main-res.syso _.rc
	$(GOPROPS) \
	GOOS=windows GOARCH=$* \
	CXX=x86_64-w64-mingw32.static-g++ CC=x86_64-w64-mingw32.static-gcc PKG_CONFIG=x86_64-w64-mingw32.static-pkg-config \
	CGO_CFLAGS_ALLOW=".*" CGO_LDFLAGS_ALLOW=".*" \
	go build -ldflags="$(LDFLAGS) -H=windows" $(GOFLAGS) -o $(BINOUT)_$*.exe

run: current
	@echo ===== RUN =====
	@./$(BINOUT)

.PHONY: all generate test current run windows_amd64 linux_386 linux_amd64 linux_arm linux_arm64 linux_mipsle
