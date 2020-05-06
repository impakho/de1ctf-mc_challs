## misc - mc_easybgm

文件列表：
[mp3.py](https://pastebin.com/8EeQ8YCz)

源自于上一年在某个比赛的启发，发现可以在mp3每一帧的保留位处隐写，当时写了实现脚本，现在发现弄丢了，然后重新写了一个。

详情可见 `mp3.py`。

[https://en.wikipedia.org/wiki/MP3#File_structure](https://en.wikipedia.org/wiki/MP3#File_structure)

`De1CTF{W31c0m3_t0_Mi73CR4Ft_W0r1D_D3jAVu!}`

## misc - mc_joinin

文件列表：
[exp.go](https://pastebin.com/Ez5khLu5)
[types.go](https://pastebin.com/rWAPYaRp)

因为想要出一道mc题，然后红石题以前比赛已经出现过，所以就没有出了。

最近在搞网络协议开发这块，所以熟悉了一下mc协议，发现可以在此处作文章。回想起中科大校赛黑曜石浏览器那题改ua版本号，发现协议通讯也涉及版本验证，所以出了这么一道题。

去重新写协议很麻烦，偷懒在某hub找到一个轮子[TyphoonMC/TyphoonLimbo](https://github.com/TyphoonMC/TyphoonLimbo)，直接魔改。

游戏开放在 `25565` 端口。

添加服务器到列表里，发现是 `MC2020` 服务端。

![](https://i.loli.net/2020/05/05/riEdb1Cu59RZcJT.png)

从官网得知 `MC2020` 是基于 `1.12` 开发的。

```
Minecraft 20.20 is developed by De1ta Team based on 1.12
```

所以我们可以用 `1.12` 的协议去通讯。

这里有两种实现方法。第一种是 MITM 中间人篡改通讯数据包里的版本号，绕过验证，成功登录游戏。第二种是直接模拟通讯协议，实现通讯。

参考资料：

```
Minecraft 1.12 Protocol(Version: 335) Wiki Page
https://wiki.vg/index.php?title=Protocol&oldid=13223
```

`exp.go` 包含两种解法。

flag藏在 `imgur` 的图片里。

![](https://i.loli.net/2020/05/05/S72oQOWhYrAcBa6.png)

`De1CTF{MC2020_Pr0to3l_Is_Funny-ISn't_It?}`

## misc - mc_champion

文件列表：
[exp.go](https://pastebin.com/964gxEzg)
[types.go](https://pastebin.com/rWAPYaRp)

此题基于轮子 [TyphoonMC/TyphoonLimbo](https://github.com/TyphoonMC/TyphoonLimbo) 魔改而来。

由于这个轮子的原因，所以数据包与市面上的MC客户端不是特别兼容，会经常掉线，尤其是网络不好的情况。所以这题建议模拟 `1.12` 通讯协议，实现通讯。

当你进入游戏，此时处于虚空时间，除了对话，其他功能都无法使用。熟悉命令行的选手，不难发现 `/help` 命令，这是一个文字游戏。

```
[ADMIN]
/help -> show the usage
/uuid -> show your uuid
/status -> show your status
/items -> show your items
/exchange -> make some exchange
/shop -> list all category
/shop [category_id] -> list items in category
/buy [item_id] -> buy the item
/use [item_id] -> use your item
/attack -> attack the BOSS
```

玩家的物品都存储在一个 `slice` 列表里，而且每一个物品都包含以下属性，fuzz一下也不难发现这点。

> Price / Attack / Shield / HP / Food / XP

漏洞函数位于 `exhcange` -> `random sell`。 这个漏洞是我平时写代码时发现的，他会触发大概像 `slice` 出栈的功能，但是由于返回值顺序的问题，导致返回了错误的值。大致代码如下：

```
func slicePop(s []int, i uint) (r []int, e int) {
    if len(s) == 0 {
        return []int{}, -1
    }else if len(s) == 1 {
        return []int{}, s[0]
    }
    if i >= uint(len(s)) {
        return s[1:], s[0]
    }
    return append(s[:i], s[i + 1:]...), s[i]
}
```

我这里的解法是，通过不断调用此功能，获得足够多的金钱（大于200），然后使用一个TNT去打败boss。

当然，赛后发现还有其他解法，只需要最终打败boss即可。

详情可见 `exp.go`。

当你打败boss，你将得到编码信息。简单地进行 `base32解码` 和 `rot13变换` ，你将得到 flag 和一个隐藏命令 `/MC2020-DEBUG-VIEW:-)`。

`De1CTF{S3cur3_UsAG3_0f-GO_Slice~}`

## reverse - mc_ticktock

文件列表：
[exp.go](https://pastebin.com/YtCdrzCd)
[types.go](https://pastebin.com/rWAPYaRp)
[crypt.go](https://pastebin.com/CYtng7mR)

上一题提到的 `/MC2020-DEBUG-VIEW:-)` 隐藏命令，是管理员用于读取指定 `uuid` 玩家 `log日志文件` 的命令。

很容易联想到 `目录穿越` 去实现 `任意文件读取`。

接下来我们可以读取以及逆向 `../../../../../../../proc/self/exe` (`go run exp.go types.go crypt.go -s1`)

文件含有符号表。

在 `main_main` 函数开头可以找到以下代码。

```
_, err := os.Stat("webserver")
if err != nil {
    log.Fatal("webserver not found")
}
```

然后我们继续读 `../../../../../../../proc/self/cwd/webserver` (`go run exp.go types.go crypt.go -s2`)

文件没有符号表，可以使用 [IDAGolangHelper](https://github.com/sibears/IDAGolangHelper) 进行恢复。

文件有以下三个功能。

1. `http://:80/ticktock?text={text}`

这里会有一个修改过的 SM4 加密算法，会对 {text} 进行加密，然后返回密文，同时与预设密文进行比较。比较相同，会将你的ip记录下来，然后你才能使用功能二和功能三，此时明文就是本题的flag值。

记录ip的表每20分钟（20分钟是mc里一昼夜的时间）清空一次。

加密过程大致如下：

```
KEY := Sha256([]byte("de1ctf-mc2020"))
NONCE := Sha256([]byte("de1ta-team"))[:24]
c, _ := crypt.NewCipher(KEY[:16])
s := cipher.NewCFBEncrypter(c, NONCE[:16])
plain := []byte("example plain text")
buff := make([]byte, len(plain))
s.XORKeyStream(buff, plain)
```

2. `http://:80/webproxy`

由于平时在开发网络协议，所以出题时也写了个代理，作为考点，想让选手访问内部网络的题。然后没想到有点难度，最后比赛过程中，这部分被临时砍掉了，出题人翻车了。

使用这个代理功能，可以实现 `http代理` 和 `端口扫描` 的功能。

剩下三道题 `mc_realworld` & `mc_logclient` & `mc_noisemap` 都需要使用这个代理去访问。

三道题的ip，可以使用 `/MC2020-DEBUG-VIEW:-) ../../../../../etc/hosts` 读取 `/etc/hosts` 来获取。

如何使用呢？发 POST 请求，POST 请求 BODY 使用 `chacha20` 加密。加密过程如下：

```
KEY := Sha256([]byte("de1ctf-mc2020"))
NONCE := Sha256([]byte("de1ta-team"))[:24]
cipher, _ := chacha20.NewUnauthenticatedCipher(KEY[:], NONCE[:])
body := []byte("127.0.0.1:80|GET /assets/ HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n")
buff := make([]byte, len(body))
cipher.XORKeyStream(buff, body)
```

3. `tcp://:8080/`

一个自定义TCP代理，用于 `mc_realworld` 的游戏服务。双向流量使用上面提到的 `chacha20` 加密。

使用命令 `go run exp.go types.go crypt.go -s3`，可以得到 flag。

`De1CTF{t1Ck-t0ck_Tlck-1ocK_MC2O20_:)SM4}`

## web - mc_logclient

文件列表：
[exp.go](https://pastebin.com/qZxVNaqm)
[types.go](https://pastebin.com/rWAPYaRp)

这是一个 `minecraft log web client`。可以用于读取所有用户的日志。

默认python3.8环境，iptables禁止对外通讯，当然白名单了icmp，可以使用ping进行外带。（赛后发现出现非预期，由于多个题目环境处于同一网络，有队伍使用 mc 协议的对话功能带出数据）

玩家的日志都存储在 logs 文件夹里，以玩家 uuid 作为文件名。logs 文件夹以 只读方式 挂载到环境里。

一个简单的 ssti，黑名单过滤了大部分关键词。在 `python 3.7` 以后，有一个新的函数 `sys.breakpointhook()` 可以通过它起一个调试器，进行任意代码执行。

赛后发现，由于黑名单不完善出现非预期，可以通过 `\x` 或者 `request.args` 进行绕过，直接进行 ssti，不需要调用 `/write` 功能。

首先我们需要在游戏对话框里，进行先 payload 的操作。开头最好加上 `/` 将 payload 作为命令的形式进行隐藏，防止向其他选手泄露 payload。

payload 如下：

```
/{{[].__class__.__base__.__subclasses__()[133].__init__.__globals__['sys']['breakpointhook']()}}
```

然后访问 `/read?work={work}&filename={uuid}` 触发 ssti。

大概有 30秒 的时间，去调用 `/write` 往 `pdb` 去写命令。

详情可见 `exp.go`。

`De1CTF{MC_L0g_C1ieNt-t0-S1mPl3_S2Tl~}`

## crypto - mc_noisemap

文件列表：
[exp.js](https://pastebin.com/P6YuDdBF)
[package.json](https://pastebin.com/gqwxugFC)
[www/map.html](https://pastebin.com/vngikVFN)
[www/assets/jquery.min.js](https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js)
[www/assets/noisemap.js](https://pastebin.com/EDna6NcL)
[www/assets/p5.dom.js](https://cdnjs.cloudflare.com/ajax/libs/p5.js/0.9.0/addons/p5.dom.js)
[www/assets/p5.js](https://cdnjs.cloudflare.com/ajax/libs/p5.js/0.9.0/p5.js)

出题时了解到 noise map 的东西，在很多游戏里都用到，包括mc。由于比较懒，查了一下，找到一个开源轮子 [erdavids/Hex-Map](https://github.com/erdavids/Hex-Map)，想想可以做个图像识别的题。

应该会有更好的解法，就不多说了，详情见 `exp.js` 吧。

`De1CTF{MCerrr_L0v3_P3r1iN_N0IsE-M4p?}`

## pwn - mc_realworld

文件列表：
[exp.py](https://pastebin.com/FrVCXW42)
[requirements.txt](https://pastebin.com/t37qDEYw)

出题人继续挖与mc相关的材料。找到了这个 [fogleman/Craft](https://github.com/fogleman/Craft) 。

游戏还蛮好的，服务端用 python 写，pwn 服务端有点难下手，而且搅屎情况更不好控制，所以打算从客户端处下手。在客户端处，埋下一个简单的 bof，在用户聊天时触发。为了防止选手集体挨打，所以限制了@特定用户时才能触发。

client to client的pwn题变得有趣多了，像极了A&D模式。

回传flag成为一个问题，中间隔着一个server。预设解法是通过聊天功能回传，只能@自己，因为公屏黑名单了关键词 `De1CTF`，防止 flag 意外泄露。

非预期回传方式，参考上面 `mc_logclient` wp 的说明。

漏洞点位于客户端 `add_messages` 函数。你可以通过 `bindiff` 找到它（需要保证编译环境一致，编译器flags一致）。代码大概如下：

```
if (text[0] == '@' && strlen(text) > 192) {
    text = text + 1;
    char *body = text + 32;
    size_t length;
    char *plain = base64_decode(body, strlen(body), &length);
    char message[16] = {0};
    memcpy(&message, plain, length);
    printf("%8s", &message);
    return;
}
```

很明显，简单 bof。

更多细节详见：exp.py。

`De1CTF{W3_L0vE_D4nge2_ReA1_W0r1d1_CrAft!2233}`