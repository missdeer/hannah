chrome.contextMenus.create({
	"title": chrome.i18n.getMessage("play_title"),
	"contexts": ["link"],
	"targetUrlPatterns":["*://music.163.com/*discover/toplist?id=*",
						"*://music.163.com/*playlist?id=*",
						"*://music.163.com/*my/m/music/playlist?id=*",
						"*://y.qq.com/n/ryqq/playlist/*",
						"*://www.kugou.com/yy/special/single/*",
						"*://*.kuwo.cn/playlist_detail/*",
						"*://music.migu.cn/v3/music/playlist/*",
					
						"*://music.163.com/*song?id=*",
						"*://y.qq.com/n/ryqq/songDetail/*",
						"*://www.kugou.com/song/*hash=*",
						"*://*.kuwo.cn/play_detail/*",
						"*://music.migu.cn/v3/music/song/*",
					
						"*://music.163.com/weapi/v1/artist/*",
						"*://music.163.com/*artist?id=*",
						"*://y.qq.com/n/ryqq/singer/*",
						"*://*.kuwo.cn/singer_detail/*",
						"*://music.migu.cn/v3/music/artist/*",

						"*://music.163.com/weapi/v1/album/*",
						"*://music.163.com/*album?id=*",
						"*://y.qq.com/n/ryqq/albumDetail/*",
						"*://*.kuwo.cn/album_detail/*",
						"*://music.migu.cn/v3/music/album/*"
					],
	"id": "hannah_link",
	"onclick": function(info, tab) {
		if (validateUrl(info.linkUrl)) {
			access("hannah://play?url=" + encodeURIComponent(info.linkUrl));
		}
	}
});

chrome.contextMenus.create({
	"title": chrome.i18n.getMessage("play_title"),
	"contexts": ["page"],
	"documentUrlPatterns":["*://music.163.com/*discover/toplist?id=*",
						"*://music.163.com/*playlist?id=*",
						"*://music.163.com/*my/m/music/playlist?id=*",
						"*://y.qq.com/n/ryqq/playlist/*",
						"*://www.kugou.com/yy/special/single/*",
						"*://*.kuwo.cn/playlist_detail/*",
						"*://music.migu.cn/v3/music/playlist/*",

						"*://music.163.com/*song?id=*",
						"*://y.qq.com/n/ryqq/songDetail/*",
						"*://www.kugou.com/song/*hash=*",
						"*://*.kuwo.cn/play_detail/*",
						"*://music.migu.cn/v3/music/song/*",

						"*://music.163.com/weapi/v1/artist/*",
						"*://music.163.com/*artist?id=*",
						"*://y.qq.com/n/ryqq/singer/*",
						"*://*.kuwo.cn/singer_detail/*",
						"*://music.migu.cn/v3/music/artist/*",

						"*://music.163.com/weapi/v1/album/*",
						"*://music.163.com/*album?id=*",
						"*://y.qq.com/n/ryqq/albumDetail/*",
						"*://*.kuwo.cn/album_detail/*",
						"*://music.migu.cn/v3/music/album/*"
					],
	"id": "hannah_play_page",
	"onclick": function(info, tab) {
		if (validateUrl(info.pageUrl)) {
			access("hannah://play?url=" + encodeURIComponent(info.pageUrl));
		}
	}	
});

chrome.contextMenus.create({
	"title": chrome.i18n.getMessage("play_title"),
	"contexts": ["audio"],
	"id": "hannah_play_audio",
	"onclick": function(info, tab) {
		access("hannah://play?url=" + encodeURIComponent(info.linkUrl));
	}
});

chrome.browserAction.onClicked.addListener(function(tab) {
	if (validateUrl(tab.url)) {
		access("hannah://play?url=" + encodeURIComponent(tab.url));
	}
});

function validateUrl(url) {
	const patterns = [
		/^https?:\/\/music\.163\.com\/(#\/)?discover\/toplist\?id=(\d+)/g,
		/^https?:\/\/music\.163\.com\/(#\/)?playlist\?id=(\d+)/g,
		/^https?:\/\/music\.163\.com\/(#\/)?my\/m\/music\/playlist\?id=(\d+)/g,
		/^https?:\/\/y\.qq\.com\/n\/ryqq\/playlist\/(\d+)/g,
		/^https?:\/\/www\.kugou\.com\/yy\/special\/single\/(\d+)\.html/g,
		/^https?:\/\/(www\.)?kuwo\.cn\/playlist_detail\/(\d+)/g,
		/^https?:\/\/music\.migu\.cn\/v3\/music\/playlist\/(\d+)/g,

		/^https?:\/\/music\.163\.com\/(#\/)?song\?id=(\d+)/g,
		/^https?:\/\/y\.qq\.com\/n\/ryqq\/songDetail\/(\w+)/g,
		/^https?:\/\/www\.kugou\.com\/song\/#hash=([0-9A-F]+)/g,
		/^https?:\/\/(www\.)?kuwo.cn\/play_detail\/(\d+)/g,
		/^https?:\/\/music\.migu\.cn\/v3\/music\/song\/(\d+)/g,

		/^https?:\/\/music\.163\.com\/weapi\/v1\/artist\/(\d+)/g,
		/^https?:\/\/music\.163\.com\/(#\/)?artist\?id=(\d+)/g,
		/^https?:\/\/y\.qq\.com\/n\/ryqq\/singer\/(\w+)/g,
		/^https?:\/\/(www\.)?kuwo\.cn\/singer_detail\/(\d+)/g,
		/^https?:\/\/music\.migu\.cn\/v3\/music\/artist\/(\d+)/g,

		/^https?:\/\/music\.163\.com\/weapi\/v1\/album\/(\d+)/g,
		/^https?:\/\/music\.163\.com\/(#\/)?album\?id=(\d+)/g,
		/^https?:\/\/y\.qq\.com\/n\/ryqq\/albumDetail\/(\w+)/g,
		/^https?:\/\/(www\.)?kuwo\.cn\/album_detail\/(\d+)/g,
		/^https?:\/\/music\.migu\.cn\/v3\/music\/album\/(\d+)/g
	];
	let valid = false;
	patterns.some(function (pattern) {
		if (pattern.test(url)) {
			valid = true;
			return true;
		}
	});
	if (!valid) {
		alert(chrome.i18n.getMessage("unsupported_link") + ": " + url);
	}
	return valid;
}

function access(url) {
	chrome.tabs.update({
		"url": url,
		"active": false
	});
}


