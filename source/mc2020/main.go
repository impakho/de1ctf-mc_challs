/**
* @Author: impakho
* @Date: 2020/04/12
* @Github: https://github.com/impakho
*/

package main

import (
	"fmt"
	t "github.com/TyphoonMC/TyphoonCore"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

const COPYRIGHT_AUTHOR = "impakho"
const COPYRIGHT_DATE = "2020/04/12"
const COPYRIGHT_GITHUB = "https://github.com/impakho"

const DEBUG = false
var core *t.Core
var players map[string]*player
var playersLock sync.Mutex

const USAGE = `
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
`

func main() {
	rand.Seed(time.Now().UnixNano())

	if !DEBUG {
		_, err := os.Stat("webserver")
		if err != nil {
			log.Fatal("webserver not found")
		}
	}

	players = make(map[string]*player, 0)

	core = t.Init()
	core.SetBrand("De1Ta He4dl3ss M1neCrAft SeRv3r")

	loadConfig(core)

	core.On(func(e *t.PlayerJoinEvent) {
		if len(e.Player.GetName()) > 32 {
			e.Player.Kick("name too long")
			return
		}

		go func() {
			f, err := os.Create(fmt.Sprintf("logs/%s", e.Player.GetUUID()))
			if err == nil {
				f.Close()
			}
		}()

		playersLock.Lock()
		players[e.Player.GetUUID()] = NewPlayer(e.Player)
		go players[e.Player.GetUUID()].Update()
		playersLock.Unlock()

		go writeLog(e.Player, "join")

		msg := t.ChatMessage("Welcome to ")
		tmsg := t.ChatMessage("[De1Ta He4dl3ss M1neCrAft Te2t SeRv3r]")
		tmsg.SetColor(&t.ChatColorPink)
		msg.SetExtra([]t.IChatComponent{
			tmsg,
			t.ChatMessage(" !"),
		})
		e.Player.SendMessage(msg)

		msg = t.ChatMessage("Hello, ")
		tmsg = t.ChatMessage(e.Player.GetName())
		tmsg.SetColor(&t.ChatColorAqua)
		msg.SetExtra([]t.IChatComponent{
			tmsg,
			t.ChatMessage(" !"),
		})
		e.Player.SendMessage(msg)

		setBossHealth(e.Player, 1.0)

		if &playerListHF != nil {
			e.Player.WritePacket(&playerListHF)
		}
	})

	core.On(func(e *t.PlayerChatEvent) {
		playerChat(e)
	})

	core.On(func(e *t.PlayerQuitEvent) {
		go func() {
			os.Remove(fmt.Sprintf("logs/%s", e.Player.GetUUID()))
		}()

		go writeLog(e.Player, "quit")

		playersLock.Lock()
		delete(players, e.Player.GetUUID())
		playersLock.Unlock()
	})

	core.Start()
}

func writeLog(player *t.Player, message string) {
	f, err := os.OpenFile("log/log.txt", os.O_WRONLY|os.O_APPEND, 0222)
	if err == nil {
		f.Write([]byte(fmt.Sprintf("%s (%s) [%s] <%s> %s\n", time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"), player.GetAddr(), player.GetUUID(), player.GetName(), message)))
		f.Close()
	}
}

func setBossHealth(player *t.Player, health float32) {
	uuid := player.GetUUID()
	playersLock.Lock()
	_, ok := players[uuid]
	if ok {
		players[uuid].BossBar.Health = health
		go player.WritePacket(&players[uuid].BossBar)
	}
	playersLock.Unlock()
}

func sendTagMsg(player *t.Player, tag string, message string) {
	msg := t.ChatMessage("")
	name := t.ChatMessage(tag)
	name.SetColor(&t.ChatColorRed)
	msg.SetExtra([]t.IChatComponent{
		t.ChatMessage("["),
		name,
		t.ChatMessage("] "),
		t.ChatMessage(message),
	})
	player.SendMessage(msg)
}

func playerChat(e *t.PlayerChatEvent) {
	playersLock.Lock()
	p, ok := players[e.Player.GetUUID()]
	if ok && p != nil {
		if !time.Now().After(time.Unix(0, p.Time).Add(1 * time.Second)) {
			go sendTagMsg(e.Player, "SYSTEM", "calm down")
			playersLock.Unlock()
			return
		}else{
			p.Time = time.Now().UnixNano()
		}
		if p.HP <= 0 {
			e.Player.Kick("You died!")
			playersLock.Unlock()
			return
		}
		if p.FOOD <= 0 {
			p.HP -= 3
			if p.HP < 0 {
				p.HP = 0
			}
			go p.Update()
		}else if p.FOOD == 20 && p.HP < 20.0 {
			p.HP += 3
			if p.HP > 20.0 {
				p.HP = 20.0
			}
			go p.Update()
		}
	}else{
		playersLock.Unlock()
		return
	}
	playersLock.Unlock()
	if len(e.Message) > 128 {
		go sendTagMsg(e.Player, "SYSTEM", "message is too long")
		return
	}
	go writeLog(e.Player, fmt.Sprintf("chat %s", e.Message))
	logMsg := e.Message[:]
	go func() {
		f, err := os.OpenFile(fmt.Sprintf("logs/%s", e.Player.GetUUID()), os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			f.Write([]byte(fmt.Sprintf("%d <%s> %s\n", time.Now().Unix(), e.Player.GetName(), logMsg)))
			f.Close()
		}
	}()
	if e.Message[:1] == "/" {
		e.Message = e.Message[1:]
		if e.Message == "help" {
			sendTagMsg(e.Player, "ADMIN", USAGE)
			return
		}
		if e.Message == "uuid" {
			sendTagMsg(e.Player, "PLAYER", e.Player.GetUUID())
			return
		}
		if e.Message == "status" {
			statusHandler(e)
			return
		}
		if e.Message == "items" {
			itemsHandler(e)
			return
		}
		if len(e.Message) >= 8 && e.Message[:8] == "exchange" {
			exchangeHandler(e)
			return
		}
		if len(e.Message) >= 3 && e.Message[:3] == "buy" {
			shopHandler(e)
			return
		}
		if len(e.Message) >= 4 && e.Message[:4] == "shop" {
			shopHandler(e)
			return
		}
		if len(e.Message) >= 3 && e.Message[:3] == "use" {
			useHandler(e)
			return
		}
		if e.Message == "attack" {
			attackHandler(e)
			return
		}
		if len(e.Message) >= 20 && e.Message[:20] == "MC2020-DEBUG-VIEW:-)" {
			debugViewHandler(e)
			return
		}
		sendTagMsg(e.Player, "SYSTEM", "invalid command~")
		return
	}

	core.GetPlayerRegistry().ForEachPlayerAsync(func(player *t.Player){
		msg := t.ChatMessage("")
		msg.SetExtra([]t.IChatComponent{
			t.ChatMessage("<"),
			t.ChatMessage(e.Player.GetName()),
			t.ChatMessage("> "),
			t.ChatMessage(e.Message),
		})
		player.SendMessage(msg)
	})
}
