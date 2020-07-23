package provider

import (
	"log"
	"testing"
)

func TestCaesar(t *testing.T) {
	location := "8%2.66E29__FD%e5d%6aeb3apt3634pd2D39%151.k%c6e528n883F%3lcx2%Eu3p2b877s_DuD6s%6%D%27E6me5a%b2F.e%%8151.ci63%rDs4d5d5_nnp11_3u2156762pyEc56a%xt22%7E6moaeD5a2i6bfa5ceus53uDt67Ef%943%fcE42i%FF5762pdmx8Et3d2fd%%ltl_9%s%iv76n589%37deaFa215E%943eip6%i2%4%152iilt52e5di%9%E722D878fsmF4%8589%%_i42o%3b5bE6ep%s26rE%d5839__6Bee1d1i885%E7233_r%6n2D6E38un%2%9ui%3%E7D%3lv69ae1"
	result := "//s128.xiami.net/868/14868/503808/1770906987_3162492_l.mp3?ccode=xiami__&expire=86400&duration=232&psid=a24624b6ebdbf0b85fd1b337da08a755&ups_client_netip=null&ups_ts=1595294613&ups_userid=0&utid=&vid=1770906987&fn=1770906987_3162492_l.mp3&vkey=B60f78e9caccd7ea60e81eeb64afd152a"
	if res, err := caesar(location); err != nil || res != result {
		t.Error(err, res)
	}
}

func TestXiami_HotPlaylist(t *testing.T) {
	p := GetProvider("xiami")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestXiami_PlaylistDetail(t *testing.T) {
	p := GetProvider("xiami")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestXiami_Search(t *testing.T) {
	p := GetProvider("xiami")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.Search("backstreet", 1, 25)
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found songs for backstreet")
	}
}

func TestXiami_SongDetail(t *testing.T) {
	p := GetProvider("xiami")
	if p == nil {
		t.Error("can't get provider")
	}

	s, err := p.SongDetail(Song{ID: "1769262490"})
	if err != nil {
		t.Error(err)
	}
	log.Println(s.URL)
}

func TestXiami_Name(t *testing.T) {
	p := GetProvider("xiami")
	if p == nil {
		t.Error("can't get provider")
	}
	if p.Name() != "xiami" {
		t.Error("provider name mismatched")
	}
}
