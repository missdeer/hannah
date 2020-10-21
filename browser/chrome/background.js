chrome.contextMenus.create({
	"title": chrome.i18n.getMessage("play_title"),
	"contexts": ["link"],
	"targetUrlPatterns":["*://*/*.torrent"],
	"id": "hannah",
	"onclick": onClickPlay
});

chrome.browserAction.onClicked.addListener(function(tab) {
	access("hannah://dialPad?address=");
});

function onClickPlay(info, tab) {
	var numberToCall = info.selectionText.replace(/\s/g, '');
	if (numberToCall.length == 0) {
		return;
	}
	access("hannah://" + numberToCall);
}

function access(url) {
	chrome.tabs.update({
		"url": url,
		"active": false
	});
}


