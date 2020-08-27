HANNAH:=hannah
RP=rp
UNAME_S:=$(shell uname -s)
ifneq ($(UNAME_S),Darwin)
	UNAME_O:=$(shell uname -o)
	ifeq ($(UNAME_O),Msys)
		HANNAH:=hannah.exe
		RP=rp.exe
	endif
endif
RPFULLPATH:=cmd/reverseProxy/$(RP)
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
    util/rand.go \
    config/config.go \
    config/cfg_darwin.go \
    config/cfg_linux.go \
    config/config_test.go \
    config/cfg_windows.go \
    config/cfg_bsd.go \
    input/source.go \
    provider/provider.go \
    provider/kuwo.go \
    provider/provider_test.go \
    provider/xiami_test.go \
    provider/migu_test.go \
    provider/xiami.go \
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
    main.go \
    rp/reverseproxy.go \
    rp/limit.go \
    rp/album.go \
    rp/artist.go \
    rp/inchina.go \
    rp/info.go \
    rp/m3ulink.go \
    rp/playlist.go \
    rp/search.go \
    rp/song.go \
    media/media.go \
    media/download.go \
    media/m3u.go

.PHONY: all
all: $(RPFULLPATH) $(HANNAH)

$(RPFULLPATH): $(CHECKS) cmd/reverseProxy/main.go
	cd cmd/reverseProxy && go build -ldflags="-s -w" -o $(RP)

$(HANNAH): $(CHECKS)
	env CGO_ENABLED=1 go build -ldflags="-s -w" -o $(HANNAH)
	if [ "$(UNAME_S)" = "Darwin" ]; then install_name_tool -change @loader_path/libbass.dylib @executable_path/output/bass/lib/darwin/amd64/libbass.dylib hannah; fi
	go mod tidy

.PHONY: clean
clean:
	env CGO_ENABLED=1 go clean
	rm -f $(HANNAH) $(RPFULLPATH)

.PHONY: test 
test:
	go test ...
