chrome.contextMenus.create({
	"title": chrome.i18n.getMessage("play_title"),
	"contexts": ["link"],
	"targetUrlPatterns":["*://music.163.com/*discover/toplist?id=*",
						"*://music.163.com/*playlist?id=*",
						"*://music.163.com/*my/m/music/playlist?id=*",
						"*://www.xiami.com/collect/*",
						"*://y.qq.com/n/yqq/playlist/*",
						"*://www.kugou.com/yy/special/single/*",
						"*://*.kuwo.cn/playlist_detail/*",
						"*://music.migu.cn/v3/music/playlist/*",
					
						"*://music.163.com/*song?id=*",
						"*://www.xiami.com/song/*",
						"*://y.qq.com/n/yqq/song/*",
						"*://www.kugou.com/song/*hash=*",
						"*://*.kuwo.cn/play_detail/*",
						"*://music.migu.cn/v3/music/song/*",
					
						"*://music.163.com/weapi/v1/artist/*",
						"*://music.163.com/*artist?id=*",
						"*://y.qq.com/n/yqq/singer/*",
						"*://www.xiami.com/artist/*",
						"*://www.xiami.com/list/scene=artist&type=*",
						"*://*.kuwo.cn/singer_detail/*",
						"*://music.migu.cn/v3/music/artist/*",

						"*://music.163.com/weapi/v1/album/*",
						"*://music.163.com/*album?id=*",
						"*://y.qq.com/n/yqq/album/*",
						"*://www.xiami.com/album/*",
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
						"*://www.xiami.com/collect/*",
						"*://y.qq.com/n/yqq/playlist/*",
						"*://www.kugou.com/yy/special/single/*",
						"*://*.kuwo.cn/playlist_detail/*",
						"*://music.migu.cn/v3/music/playlist/*",

						"*://music.163.com/*song?id=*",
						"*://www.xiami.com/song/*",
						"*://y.qq.com/n/yqq/song/*",
						"*://www.kugou.com/song/*hash=*",
						"*://*.kuwo.cn/play_detail/*",
						"*://music.migu.cn/v3/music/song/*",

						"*://music.163.com/weapi/v1/artist/*",
						"*://music.163.com/*artist?id=*",
						"*://y.qq.com/n/yqq/singer/*",
						"*://www.xiami.com/artist/*",
						"*://www.xiami.com/list/scene=artist&type=*",
						"*://*.kuwo.cn/singer_detail/*",
						"*://music.migu.cn/v3/music/artist/*",

						"*://music.163.com/weapi/v1/album/*",
						"*://music.163.com/*album?id=*",
						"*://y.qq.com/n/yqq/album/*",
						"*://www.xiami.com/album/*",
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
		/^https?:\/\/www\.xiami\.com\/collect\/(\d+)/g,
		/^https?:\/\/y\.qq\.com\/n\/yqq\/playlist\/(\d+)\.html/g,
		/^https?:\/\/www\.kugou\.com\/yy\/special\/single\/(\d+)\.html/g,
		/^https?:\/\/(www\.)?kuwo\.cn\/playlist_detail\/(\d+)/g,
		/^https?:\/\/music\.migu\.cn\/v3\/music\/playlist\/(\d+)/g,

		/^https?:\/\/music\.163\.com\/(#\/)?song\?id=(\d+)/g,
		/^https?:\/\/www\.xiami\.com\/song\/(\w+)/g,
		/^https?:\/\/y\.qq\.com\/n\/yqq\/song\/(\w+)\.html/g,
		/^https?:\/\/www\.kugou\.com\/song\/#hash=([0-9A-F]+)/g,
		/^https?:\/\/(www\.)?kuwo.cn\/play_detail\/(\d+)/g,
		/^https?:\/\/music\.migu\.cn\/v3\/music\/song\/(\d+)/g,

		/^https?:\/\/music\.163\.com\/weapi\/v1\/artist\/(\d+)/g,
		/^https?:\/\/music\.163\.com\/(#\/)?artist\?id=(\d+)/g,
		/^https?:\/\/y\.qq\.com\/n\/yqq\/singer\/(\w+)\.html/g,
		/^https?:\/\/www\.xiami\.com\/artist\/(\w+)/g,
		/^https?:\/\/www\.xiami\.com\/list\?scene=artist&type=\w+&query={%22artistId%22:%22(\d+)%22}/g,
		/^https?:\/\/www\.xiami\.com\/list\?scene=artist&type=\w+&query={"artistId":"(\d+)"}/g,
		/^https?:\/\/(www\.)?kuwo\.cn\/singer_detail\/(\d+)/g,
		/^https?:\/\/music\.migu\.cn\/v3\/music\/artist\/(\d+)/g,

		/^https?:\/\/music\.163\.com\/weapi\/v1\/album\/(\d+)/g,
		/^https?:\/\/music\.163\.com\/(#\/)?album\?id=(\d+)/g,
		/^https?:\/\/y\.qq\.com\/n\/yqq\/album\/(\w+)\.html/g,
		/^https?:\/\/www\.xiami\.com\/album\/(\w+)/g,
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


