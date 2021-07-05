GITCOMMIT:=$(shell git describe --dirty --always)
HANNAH:=hannah
RP=rp
LIBRP=librp.a
LIBBUILDOPTS:=-v -buildmode=c-archive
UNAME_S:=$(shell uname -s)
UNAME_M:=$(shell uname -m)
ifneq ($(UNAME_M),arm64)
	GOARCH:=amd64
else
	GOARCH:=arm64
endif
ifneq ($(UNAME_S),Darwin)
	UNAME_O:=$(shell uname -o)
	ifeq ($(UNAME_O),Msys)
		HANNAH:=hannah.exe
		RP=rp.exe
	endif
endif
RPFULLPATH:=cmd/reverseProxy/$(RP)
HANNAHFULLPATH:=cmd/hannah/$(HANNAH)
LIBRPFULLPATH:=lib/reverseProxy/$(LIBRP)
CHECKS:=go.mod go.sum \
	util/cryptography/ecb.go \
    util/cryptography/aes.go \
    util/cryptography/rsa.go \
    util/cryptography/des.go \
    util/conv_test.go \
    util/http.go \
    util/player.go \
    util/conv.go \
    util/json.go \
    util/linkpattern.go \
    util/rand.go \
    config/config.go \
    config/cfg_darwin.go \
    config/cfg_linux.go \
    config/config_test.go \
    config/cfg_windows.go \
    config/cfg_bsd.go \
    input/source.go \
    lyric/lrc.go \
    lyric/xtrc.go \
    provider/provider.go \
    provider/kuwo.go \
    provider/provider_test.go \
    provider/migu_test.go \
    provider/qq.go \
    provider/kuwo_test.go \
    provider/migu.go \
    provider/bilibili.go \
    provider/bilibili_test.go \
    provider/netease.go \
    provider/netease_test.go \
    provider/kugou.go \
    provider/kugou_test.go \
    provider/qq_test.go \
    output/panel.go \
    output/speaker.go \
    output/bass/bassplugins.go \
    output/bass/bass.go \
    output/bass/format.go \
    output/bass/bridge.go \
    output/bass/speaker.go \
    action/action.go \
    action/play.go \
    action/action_test.go \
    action/search.go \
    action/hot.go \
    action/playlist.go \
    cache/redis.go \
    rp/reverseproxy.go \
    rp/limit.go \
    rp/album.go \
    rp/artist.go \
    rp/inchina.go \
    rp/info.go \
    rp/inlan.go \
    rp/m3ulink.go \
    rp/playlist.go \
    rp/search.go \
    rp/song.go \
    media/media.go \
    media/download.go \
    media/m3u.go

.PHONY: all
all: $(LIBRPFULLPATH) $(RPFULLPATH) $(HANNAHFULLPATH) 

$(RPFULLPATH): $(CHECKS) cmd/reverseProxy/main.go
	cd cmd/reverseProxy && go build -ldflags="-s -w -X main.GitCommit=$(GITCOMMIT)" -o $(RP)
	if [ "$(UNAME_S)" = "Darwin" -a "$(UNAME_M)" = "arm64" ]; then cd cmd/reverseProxy && env CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-s -w -X main.GitCommit=$(GITCOMMIT)" -o $(RP).amd64; fi
	if [ "$(UNAME_S)" = "Darwin" -a "$(UNAME_M)" = "arm64" ]; then cd cmd/reverseProxy && lipo -create -output $(RP) $(RP) $(RP).amd64; rm $(RP).amd64; fi

$(HANNAHFULLPATH): $(CHECKS) cmd/hannah/main.go
	cd cmd/hannah && env CGO_ENABLED=1 go build -ldflags="-s -w -X main.GitCommit=$(GITCOMMIT)" -o $(HANNAH)
	if [ "$(UNAME_S)" = "Darwin" ]; then install_name_tool -change @loader_path/libbass.dylib @executable_path/../../output/bass/lib/darwin/amd64/libbass.dylib cmd/hannah/$(HANNAH); fi
	if [ "$(UNAME_S)" = "Darwin" -a "$(UNAME_M)" = "arm64" ]; then cd cmd/hannah && env CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-s -w -X main.GitCommit=$(GITCOMMIT)" -o $(HANNAH).amd64; fi
	if [ "$(UNAME_S)" = "Darwin" -a "$(UNAME_M)" = "arm64" ]; then install_name_tool -change @loader_path/libbass.dylib @executable_path/../../output/bass/lib/darwin/amd64/libbass.dylib cmd/hannah/$(HANNAH).amd64; fi
	if [ "$(UNAME_S)" = "Darwin" -a "$(UNAME_M)" = "arm64" ]; then cd cmd/hannah && lipo -create -output $(HANNAH) $(HANNAH) $(HANNAH).amd64; rm $(HANNAH).amd64; fi

$(LIBRPFULLPATH): $(CHECKS) lib/reverseProxy/main.go
	cd lib/reverseProxy && CGO_ENABLED=1 go build $(LIBBUILDOPTS) -ldflags="-s -w -X main.GitCommit=$(GITCOMMIT)" -o $(LIBRP)
	if [ "$(UNAME_S)" = "Darwin" -a "$(UNAME_M)" = "arm64" ]; then cd lib/reverseProxy && env CGO_ENABLED=1 GOARCH=amd64 go build $(LIBBUILDOPTS) -ldflags="-s -w -X main.GitCommit=$(GITCOMMIT)" -o $(LIBRP).amd64.a; fi
	if [ "$(UNAME_S)" = "Darwin" -a "$(UNAME_M)" = "arm64" ]; then cd lib/reverseProxy && lipo -create -output $(LIBRP) $(LIBRP) $(LIBRP).amd64.a; rm $(LIBRP).amd64.a $(LIBRP).amd64.h; fi

.PHONY: clean
clean:
	cd cmd/reverseProxy && env CGO_ENABLED=1 go clean
	cd cmd/hannah && env CGO_ENABLED=1 go clean
	cd lib/reverseProxy && env CGO_ENABLED=1 go clean
	rm -f $(HANNAHFULLPATH) $(RPFULLPATH) $(LIBRPFULLPATH) $(HANNAHFULLPATH).amd64 $(RPFULLPATH).amd64 $(LIBRPFULLPATH).amd64.a

.PHONY: test 
test:
	go test ...
