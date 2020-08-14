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
```

* use [foobar2000](http://www.foobar2000.org/) to open `old.m3u` and enjoy it.

# 在 Windows 上编译 hannah 的方法
clone 源项目代码后，直接go build . 失败。
下面就是按照报错信息提示一步步来的，所以其实编译本身并不难，就是花点时间找齐需要的东西。

为了节省找东西的时间，将步骤记录如下：

* 我是在 windows 10 2004 中，在 pwsh 7.3 环境中进行的
* 需要 GCC。因为是windows 10 64位，安装的golang 也是64位版本，所以我使用的是 msys2 的 x86_64-w64-mingw32。
* 需要自行编译 [faad2](https://github.com/knik0/faad2) 项目提供msvc sln 文件，可以用visual studio 编译。需要配置编译x64版本。
* 需要 [mpg123](mpg123.org) 直接下载x64版本即可。
* 需要 pkg-config，我选择 mysy2 安装：```pacman -S pkg-config``` 即可。

准备 libmpg123.pc 和 faad2.pc 文件，格式抄[维基百科](https://zh.wikipedia.org/wiki/Pkg-config)的就行。就是注意路径配置正确。

下面是我本地的配置：

faad2.pc:
```
prefix="E:\Downloads\faad2-2_9_2"
exec_prefix="${prefix}"
libdir="${exec_prefix}\project\msvc\x64\Release"
includedir="${exec_prefix}\include"

Name: libfaad2
Description: Loads and play aac files
Version: 2.9.2
Libs: -L${libdir} -lfaad2
Cflags: -I${includedir}
```
libmpg123.pc:
```
prefix="E:\Downloads\mpg123-1.26.3-x86-64"
exec_prefix=${prefix}
libdir=${exec_prefix}
includedir=${exec_prefix}

Name: libmpg123
Description: Loads and play mp3 files
Version: 1.26.3
Libs: -L${libdir} -lmpg123-0
Cflags: -I${includedir}
```

配置完后，使用
```pkg-config --cflags  -- faad2``` 和 ```pkg-config --cflags  -- libmpg123```
来检查配置的路径是否正确。

配置一下 PATH 环境变量，使能找到 msys2 的 gcc 和 pkg-config
```$env:PATH=$env:PATH+";C:\msys64\mingw64\bin;C:\msys64\usr\bin"```

配置一下 PKG_CONFIG_PATH，使能找到刚配置的两个pc 文件。我的两个pc 文件就放在 hannah 项目目录下，所以：
```$env:PKG_CONFIG_PATH="."``` 就行。

然后就可以在项目目录下执行 go build . 了。

会在项目目录下得到hannah.exe

但此时在命令行（pwsh）中执行，会报错，提示缺少一些 BASS 系列的库

因此还需要下载 [BASS](www.un4seen.com) 系列 dll：
* bass.dll
* bassasio.dll
* bassmix.dll
* basswasapi.dll
* 为播放flac，还需要 basswasapi.dll

至此可以使用 hannah 听歌了