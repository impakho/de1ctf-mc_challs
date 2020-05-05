package main

import (
	"fmt"
	t "github.com/TyphoonMC/TyphoonCore"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type category struct {
	Cid int
	Name string
	EnName string
}

type item struct {
	Cid int
	Id int
	Name string
	EnName string
	Price int
	Attack int
	Shield int
	HP float32
	Food int
	XP int
}

var categorys = []category{
	{Cid: 1, Name: "武器", EnName: "Weapon"},
	{Cid: 2, Name: "防具", EnName: "Armor"},
	{Cid: 3, Name: "工具", EnName: "Tool"},
	{Cid: 4, Name: "食物", EnName: "Food"},
}

var items = []item{
	{Cid: 1, Id: 1, Name: "木锄", EnName: "Wooden Hoe", Price: 7, Attack: 3, Shield: 0, HP: 0, Food: -1, XP: 4},
	{Cid: 1, Id: 2, Name: "石锹", EnName: "Stone Shovel", Price: 10, Attack: 5, Shield: 0, HP: 0, Food: -1, XP: 6},
	{Cid: 1, Id: 3, Name: "铁镐", EnName: "Iron Pickaxe", Price: 13, Attack: 8, Shield: 0, HP: 0, Food: -1, XP: 8},
	{Cid: 1, Id: 4, Name: "金斧", EnName: "Golden Axe", Price: 33, Attack: 30, Shield: 0, HP: 0, Food: -2, XP: 12},
	{Cid: 1, Id: 5, Name: "钻石剑", EnName: "Diamond Sword", Price: 56, Attack: 50, Shield: 0, HP: 0, Food: -3, XP: 20},
	{Cid: 1, Id: 6, Name: "弓", EnName: "Bow", Price: 24, Attack: 0, Shield: 0, HP: 0, Food: -1, XP: 5},
	{Cid: 2, Id: 7, Name: "铁头盔", EnName: "Golden Helmet", Price: 20, Attack: 0, Shield: 5, HP: 0, Food: 0, XP: 0},
	{Cid: 2, Id: 8, Name: "金胸甲", EnName: "Golden Chestplate", Price: 33, Attack: 0, Shield: 14, HP: 0, Food: 0, XP: 0},
	{Cid: 2, Id: 9, Name: "钻石头盔", EnName: "Diamond Helmet", Price: 50, Attack: 0, Shield: 28, HP: 0, Food: 0, XP: 0},
	{Cid: 2, Id: 10, Name: "钻石胸甲", EnName: "Diamond Chestplate", Price: 47, Attack: 27, Shield: 25, HP: 0, Food: 0, XP: 0},
	{Cid: 2, Id: 11, Name: "钻石护腿", EnName: "Diamond Leggings", Price: 42, Attack: 0, Shield: 22, HP: 0, Food: 0, XP: 0},
	{Cid: 2, Id: 12, Name: "钻石靴子", EnName: "Diamond Boots", Price: 25, Attack: 0, Shield: 15, HP: 0, Food: 0, XP: 0},
	{Cid: 3, Id: 13, Name: "雪球", EnName: "Snowball", Price: 1, Attack: 1, Shield: 0, HP: 0, Food: 0, XP: 0},
	{Cid: 3, Id: 14, Name: "炸药", EnName: "TNT", Price: 200, Attack: 1000, Shield: 0, HP: -10.2, Food: -2, XP: 0},
	{Cid: 3, Id: 15, Name: "箭", EnName: "Arrow", Price: 1, Attack: 7, Shield: 0, HP: 0, Food: 0, XP: 0},
	{Cid: 3, Id: 16, Name: "打火石", EnName: "Flint and Steel", Price: 10, Attack: 0, Shield: 0, HP: 0, Food: 0, XP: 30},
	{Cid: 3, Id: 17, Name: "金苹果", EnName: "Golden Apple", Price: 233, Attack: 0, Shield: 0, HP: 20, Food: 20, XP: 2000},
	{Cid: 3, Id: 18, Name: "闪烁的西瓜片", EnName: "Glistering Melon Slice", Price: 99, Attack: 0, Shield: 0, HP: 15, Food: 20, XP: 300},
	{Cid: 4, Id: 19, Name: "面包", EnName: "Bread", Price: 4, Attack: 0, Shield: 0, HP: 2, Food: 2, XP: 0},
	{Cid: 4, Id: 20, Name: "胡萝卜", EnName: "Carrot", Price: 3, Attack: 0, Shield: 0, HP: 1, Food: 2, XP: 0},
	{Cid: 4, Id: 21, Name: "蜘蛛眼", EnName: "Spider Eye", Price: 2, Attack: 0, Shield: 0, HP: -5, Food: 1, XP: 0},
	{Cid: 4, Id: 22, Name: "生鸡肉", EnName: "Raw Chicken", Price: 2, Attack: 0, Shield: 0, HP: -2, Food: 2, XP: 0},
	{Cid: 4, Id: 23, Name: "毒马铃薯", EnName: "Poisonous Potato", Price: 2, Attack: 0, Shield: 0, HP: -4.4, Food: 1, XP: 0},
	{Cid: 4, Id: 24, Name: "熔岩桶", EnName: "Lava Bucket", Price: 50, Attack: 0, Shield: 0, HP: -20, Food: 0, XP: 1000},
}

func shopHandler(e *t.PlayerChatEvent) {
	if e.Message == "shop" {
		output := "\n"
		for _, k := range categorys {
			output += fmt.Sprintf("[%d] %s %s\n", k.Cid, k.Name, k.EnName)
		}
		sendTagMsg(e.Player, "SHOP", output)
	}else if len(e.Message) >= 4 && e.Message[:4] == "shop" {
		sep := strings.Split(e.Message, " ")
		if len(sep) != 2 {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		if sep[0] != "shop" {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		n, err := strconv.Atoi(sep[1])
		if err != nil {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		output := "\n"
		for _, k := range items {
			if k.Cid == n {
				output += fmt.Sprintf("[%d] %s %s $%d\n", k.Id, k.Name, k.EnName, k.Price)
			}
		}
		if output == "\n" {
			sendTagMsg(e.Player, "SHOP", "\ninvalid category")
		}else{
			sendTagMsg(e.Player, "SHOP", output)
		}
	}else if e.Message[:3] == "buy" {
		sep := strings.Split(e.Message, " ")
		if len(sep) != 2 {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		if sep[0] != "buy" {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		n, err := strconv.Atoi(sep[1])
		if err != nil {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		sitem := item{}
		for _, k := range items {
			if k.Id == n {
				sitem = k
			}
		}
		if sitem == (item{}) {
			sendTagMsg(e.Player, "SHOP", "\ninvalid item")
			return
		}
		uuid := e.Player.GetUUID()
		playersLock.Lock()
		p, ok := players[uuid]
		playersLock.Unlock()
		if p == nil || !ok {
			return
		}
		if p.MONEY < sitem.Price {
			sendTagMsg(e.Player, "SHOP", "\nlack of money")
			return
		}
		playersLock.Lock()
		_, ok = players[uuid]
		if !ok {
			playersLock.Unlock()
			return
		}
		players[uuid].MONEY -= sitem.Price
		players[uuid].Items = append(players[uuid].Items, sitem.Id)
		playersLock.Unlock()
		sendTagMsg(e.Player, "SHOP", "\nbuy succ")
	}
}

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

func exchangeHandler(e *t.PlayerChatEvent) {
	if e.Message == "exchange" {
		sendTagMsg(e.Player, "EXCHANGE", "\n[1] $1 <- 10 XP\n[2] 3 XP <- $1\n[3] $1 <- 5 HUNGER\n[4] 1 HUNGER <- $3\n[5] 1 HP <- $10\n[6] RANDOM SELL\n    Sell the item at 10% off from your item list randomly.")
		return
	}else if e.Message[:8] == "exchange" {
		sep := strings.Split(e.Message, " ")
		if len(sep) != 2 {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		if sep[0] != "exchange" {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		n, err := strconv.Atoi(sep[1])
		if err != nil {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		uuid := e.Player.GetUUID()
		playersLock.Lock()
		p, ok := players[uuid]
		playersLock.Unlock()
		if p == nil || !ok {
			return
		}
		if n == 1 {
			if p.XP < 10 {
				sendTagMsg(e.Player, "EXCHANGE", "\nlack of xp")
				return
			}
			playersLock.Lock()
			players[uuid].XP -= 10
			if players[uuid].XP < 0 {
				players[uuid].XP = 0
			}
			players[uuid].MONEY += 1
			go players[uuid].Update()
			playersLock.Unlock()
			sendTagMsg(e.Player, "EXCHANGE", "\nexchange succ, get $1")
			go statusHandler(e)
		}else if n == 2 {
			if p.MONEY < 1 {
				sendTagMsg(e.Player, "EXCHANGE", "\nlack of money")
				return
			}
			playersLock.Lock()
			players[uuid].MONEY -= 1
			if players[uuid].MONEY < 0 {
				players[uuid].MONEY = 0
			}
			players[uuid].XP += 3
			go players[uuid].Update()
			playersLock.Unlock()
			sendTagMsg(e.Player, "EXCHANGE", "\nexchange succ, get 3 XP")
			go statusHandler(e)
		}else if n == 3 {
			if p.FOOD < 5 {
				sendTagMsg(e.Player, "EXCHANGE", "\nlack of food")
				return
			}
			playersLock.Lock()
			players[uuid].FOOD -= 5
			if players[uuid].FOOD < 0 {
				players[uuid].FOOD = 0
			}
			players[uuid].MONEY += 1
			go players[uuid].Update()
			playersLock.Unlock()
			sendTagMsg(e.Player, "EXCHANGE", "\nexchange succ, get $1")
			go statusHandler(e)
		}else if n == 4 {
			if p.MONEY < 3 {
				sendTagMsg(e.Player, "EXCHANGE", "\nlack of money")
				return
			}
			playersLock.Lock()
			players[uuid].MONEY -= 3
			if players[uuid].MONEY < 0 {
				players[uuid].MONEY = 0
			}
			players[uuid].FOOD += 1
			go players[uuid].Update()
			playersLock.Unlock()
			sendTagMsg(e.Player, "EXCHANGE", "\nexchange succ, get 1 HUNGER")
			go statusHandler(e)
		}else if n == 5 {
			if p.MONEY < 10 {
				sendTagMsg(e.Player, "EXCHANGE", "\nlack of money")
				return
			}
			playersLock.Lock()
			players[uuid].MONEY -= 10
			if players[uuid].MONEY < 0 {
				players[uuid].MONEY = 0
			}
			players[uuid].HP += 1
			go players[uuid].Update()
			playersLock.Unlock()
			sendTagMsg(e.Player, "EXCHANGE", "\nexchange succ, get 1 HP")
			go statusHandler(e)
		}else if n == 6 {
			playersLock.Lock()
			if len(players[uuid].Items) <= 0 {
				playersLock.Unlock()
				sendTagMsg(e.Player, "EXCHANGE", "\nlack of item")
				return
			}
			var item int
			players[uuid].Items, item = slicePop(players[uuid].Items, uint(rand.Intn(len(players[uuid].Items))))
			if item == -1 {
				playersLock.Unlock()
				sendTagMsg(e.Player, "EXCHANGE", "\nlack of item")
				return
			}
			earn := 0
			for _, h := range items {
				if h.Id == item {
					earn = int(math.Floor(float64(h.Price) * 0.9))
					break
				}
			}
			players[uuid].MONEY += earn
			itemlist := players[uuid].Items
			playersLock.Unlock()
			itemlistd := "["
			found := false
			for _, k := range itemlist {
				for _, h := range items {
					if k == h.Id {
						found = true
						itemlistd += fmt.Sprintf("\"[%d] %s %s\", ", h.Id, h.Name, h.EnName)
						break
					}
				}
			}
			if found {
				itemlistd = itemlistd[:len(itemlistd) - 2]
			}
			itemlistd += "]"
			sendTagMsg(e.Player, "EXCHANGE", fmt.Sprintf("\nexchange succ, get $%d\nyour full item list:\n%s", earn, itemlistd))
			go statusHandler(e)
		}
	}
}