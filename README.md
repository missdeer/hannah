# Hannah

Listen music

## Usage

### Show the help message
```bash
./hannah -h
```

### Play mp3 files in order
```bash
./hannah music1.mp3 music2.mp3 music3.mp3
```

### Search and play all supported music file in the directory and all sub-directories
```bash
./hannah music-directory
```

### Play mp3 files again and again in shuffle order
```bash
./hannah --repeat --shuffle music1.mp3 music2.mp3 music3.mp3
```

### Play songs by [foobar2000](http://www.foobar2000.org/) or other your favourite media player
 
* generate m3u file that includes all songs you want to play
```bash
# save all songs in the first page of the specified playlist to the specified m3u file
./hannah -a playlist-save -p qq --reverse-proxy-enabled --reverse-proxy 127.0.0.1:8888 --m3u old.m3u 7602926765

# search the keyword and save all songs in the first page of the search result to the specified m3u file
./hannah -a search-save -p qq --reverse-proxy-enabled --reverse-proxy 127.0.0.1:8888 --m3u westlife.m3u westlife
```

* launch reverse proxy
    
```bash
# normal case
./cmd/reverseProxy/rp -b 127.0.0.1:8888

# if you are in China, don't need a proxy, use `redirect` mode to improve performance
./cmd/reverseProxy/rp -b 127.0.0.1:8888 --redirect

# if you are NOT in China, need to use a proxy to access those song services
./cmd/reverseProxy/rp -b 127.0.0.1:8888 --socks5 127.0.0.1:8080

# if you have two or more network connections and want to specify one to be used by Hannah
./cmd/reverseProxy/rp -b 127.0.0.1:8888 --network-interface en1
# or
./cmd/reverseProxy/rp -b 127.0.0.1:8888 --network-interface 192.168.2.100
```

* use [foobar2000](http://www.foobar2000.org/) to open `old.m3u` and enjoy it.

### Used within web browser

* launch reverse proxy
    
```bash
# normal case
./cmd/reverseProxy/rp -b 127.0.0.1:8888

# if you are in China, don't need a proxy, use `redirect` mode to improve performance
./cmd/reverseProxy/rp -b 127.0.0.1:8888 --redirect

# if you are NOT in China, need to use a proxy to access those song services
./cmd/reverseProxy/rp -b 127.0.0.1:8888 --socks5 127.0.0.1:8080

# if you have two or more network connections and want to specify one to be used by Hannah
./cmd/reverseProxy/rp -b 127.0.0.1:8888 --network-interface en1
# or
./cmd/reverseProxy/rp -b 127.0.0.1:8888 --network-interface 192.168.2.100
```

* add a bookmark button on your web browser's Bookmarks Toolbar

Name field: `Generate M3U` or any other name you like 

Location field:

```javascript
javascript:(function(){var%20u=window.location.href;window.open("http://127.0.0.1:8888/m3u/generate?u="+encodeURIComponent(u),"_blank")})();void%200;
```

* open a playlist page within your web browser, for example:

```txt
https://music.163.com/#/playlist?id=5149363884
```

* click `Generate M3U` button created in above step
* choose your favourite media player application, such as [foobar2000](http://www.foobar2000.org/), and click `Open` button on the web browser's opening dialog, enjoy it.