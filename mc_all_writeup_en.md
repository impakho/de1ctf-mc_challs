## misc - mc_easybgm

Files:
[mp3.py](https://pastebin.com/8EeQ8YCz)

Unused bit stego. The script: `mp3.py`.

[https://en.wikipedia.org/wiki/MP3#File_structure](https://en.wikipedia.org/wiki/MP3#File_structure)

`De1CTF{W31c0m3_t0_Mi73CR4Ft_W0r1D_D3jAVu!}`

## misc - mc_joinin

Files:
[exp.go](https://pastebin.com/Ez5khLu5)
[types.go](https://pastebin.com/rWAPYaRp)

The minecraft game service is opened on default port `25565`.

We add the server to mutil-player server list. It seems that it's not supported by our client. The server is using 'MC2020' as service and our client is seems to be outdated.

![](https://i.loli.net/2020/05/05/riEdb1Cu59RZcJT.png)

From the website we know that, `MC2020` is developed based on `1.12`.

```
Minecraft 20.20 is developed by De1ta Team based on 1.12
```

Of course there is a thick. The server is still using `1.12` protocol to communicate. But in the 'protocol version checking' procedure, it's just failed.

It comes out two solutions. One, we just simply replace the 'version' in the procedure during communication between the server and the client by a custom proxy. Another, we simulate the client's function, login to the game.

Reference:

```
Minecraft 1.12 Protocol(Version: 335) Wiki Page
https://wiki.vg/index.php?title=Protocol&oldid=13223
```

`exp.go` is the exploit written in golang, containing with two solutions mentioned above.

The flag is hidden in the image from `imgur`.

![](https://i.loli.net/2020/05/05/S72oQOWhYrAcBa6.png)

`De1CTF{MC2020_Pr0to3l_Is_Funny-ISn't_It?}`

## misc - mc_champion

Files:
[exp.go](https://pastebin.com/964gxEzg)
[types.go](https://pastebin.com/rWAPYaRp)

The Minecraft game is modified basing on [TyphoonMC/TyphoonLimbo](https://github.com/TyphoonMC/TyphoonLimbo).

When you login to the game, you are in the Limbo World. Nothing you can do except chatting with other players. After fuzzing, it's not difficult to find out the `/help` command.

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

Player's items are stored in a golang slice. Each item has five attributes listed below. After a fuzzing, you may figure it out.

> Price / Attack / Shield / HP / Food / XP

The vulnerable function is located in `exhcange` -> `random sell`. It triggles a `pop from slice` liked function which return result in bad sequence caused wrong item pop out.

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

Our goal is earning enough money (more than $200) and using a TNT to defeat the boss.

Want more details? Have a look on the exploit `exp.go` written in golang.

After won the game, we got the `Encoded Message`. Simply have a `Base32 Decode` and a `rot13`. We got Flag Two of this challenge and a hidden command `/MC2020-DEBUG-VIEW:-)`. Next challenge, we will use this command.

`De1CTF{S3cur3_UsAG3_0f-GO_Slice~}`

## reverse - mc_ticktock

Files:
[exp.go](https://pastebin.com/YtCdrzCd)
[types.go](https://pastebin.com/rWAPYaRp)
[crypt.go](https://pastebin.com/CYtng7mR)

Previously on the challenge, we got a hidden command `/MC2020-DEBUG-VIEW:-)`. We could read player's log file by their's UUID. Of course, it's a classical directory traversal attack here to read any file on the challenge environment.

Let's read the service binary file `../../../../../../../proc/self/exe`, and reverse it. (`go run exp.go types.go crypt.go -s1`)

In `main_main` function, there is some code like:

```
_, err := os.Stat("webserver")
if err != nil {
    log.Fatal("webserver not found")
}
```

Go on, and read the web service binary file `../../../../../../../proc/self/cwd/webserver`, and reverse it. (`go run exp.go types.go crypt.go -s2`)

You will found three hidden functions.

1. `http://:80/ticktock?text={text}`

It will have a Modified-SM4 encryption of the {text}, and compare the cipher-text with the prefix one. If they match, you will have 20 minutes (One day-night cycle in Minecraft World) to access function two and three. In the meantime, the plain text contains the flag of this challenge.

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

It's a custom proxy service. You can use it to make HTTP request or do a TCP scanning. Remaining three challenges `mc_realworld` & `mc_logclient` & `mc_noisemap` need this proxy to access the web service.

How to use? Make a POST request to this URL. The POST body should be encrypted using `chacha20` cipher.

```
KEY := Sha256([]byte("de1ctf-mc2020"))
NONCE := Sha256([]byte("de1ta-team"))[:24]
cipher, _ := chacha20.NewUnauthenticatedCipher(KEY[:], NONCE[:])
body := []byte("127.0.0.1:80|GET /assets/ HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n")
buff := make([]byte, len(body))
cipher.XORKeyStream(buff, body)
```

3. `tcp://:8080/`

It's a custom TCP proxy service to access the game service of challenge `mc_realworld`. Inbound traffic and outbound traffic are both using `chacha20` cipher to encrypt which mentioned above.

To get the flag, try `go run exp.go types.go crypt.go -s3`.

`De1CTF{t1Ck-t0ck_Tlck-1ocK_MC2O20_:)SM4}`

## web - mc_logclient

Files:
[exp.go](https://pastebin.com/qZxVNaqm)
[types.go](https://pastebin.com/rWAPYaRp)

It's a `minecraft log web client`.

The `chatting content` of players are stored in `logs/` with their's UUID. The `logs/` is mounted as read-only directory.

A simple `SSTI` with `render_template_string` of `flask`. Almost words are blacklisted. After `python 3.7`, there is a new function `sys.breakpointhook()`, using it to run arbitrary code.

First, saying `payload` showed below in the chat box of minecraft. (Message with a '/' prefix to act as a command, for hiding your `payload` from other players)

The environment is `python 3.8`.

```
/{{[].__class__.__base__.__subclasses__()[133].__init__.__globals__['sys']['breakpointhook']()}}
```

Then, triggle `render_template_string` by visiting `/read?work={work}&filename={uuid}`.

You have 30 seconds to access `/write` for writing the `command` to `pdb`.

More details, have a check on exploit file `exp.go`.

`De1CTF{MC_L0g_C1ieNt-t0-S1mPl3_S2Tl~}`

## crypto - mc_noisemap

Files:
[exp.js](https://pastebin.com/P6YuDdBF)
[package.json](https://pastebin.com/gqwxugFC)
[www/map.html](https://pastebin.com/vngikVFN)
[www/assets/jquery.min.js](https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js)
[www/assets/noisemap.js](https://pastebin.com/EDna6NcL)
[www/assets/p5.dom.js](https://cdnjs.cloudflare.com/ajax/libs/p5.js/0.9.0/addons/p5.dom.js)
[www/assets/p5.js](https://cdnjs.cloudflare.com/ajax/libs/p5.js/0.9.0/p5.js)

The challenge is about `Image Identification`, modified basing on [erdavids/Hex-Map](https://github.com/erdavids/Hex-Map).

My solution is not perfect. Perhaps there are some good ways to solve the challenge, I believe.

Just have a check on the exploit file `exp.js`.

`De1CTF{MCerrr_L0v3_P3r1iN_N0IsE-M4p?}`

## pwn - mc_realworld

Files:
[exp.py](https://pastebin.com/FrVCXW42)
[requirements.txt](https://pastebin.com/t37qDEYw)

The challenge was modified based on a minecraft-liked game written in C [fogleman/Craft](https://github.com/fogleman/Craft).

The vulnerability function is `add_messages` located in the client binary. You can use `bindiff` to find it. In the function, there is some codes like:

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

Obviously, an easy stack BOF! Let's use `checksec` to have a look at the protection.

```
Arch:     amd64-64-little
RELRO:    Partial RELRO
Stack:    No canary found
NX:       NX enabled
PIE:      No PIE (0x400000)
```

Okay... A check-in challenge, that's what I'm thinking about.

`add_messages` will be triggled before some messages show on the console in game. Therefore, we can exploit a player's machine by @someone in the chat box. Make sure the length of text message above 192 bytes, also using `base64` to encode.

One more problem, how to get the `flag` from the victim(bot)'s machine. After digging in the binary, I find `client_talk` function. Using it to @attacker(me) follow with the `flag`, we can receive the flag at the client side.

For more details, check the expolit `exp.py`.

`De1CTF{W3_L0vE_D4nge2_ReA1_W0r1d1_CrAft!2233}`