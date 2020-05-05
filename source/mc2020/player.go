package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	t "github.com/TyphoonMC/TyphoonCore"
	uuid "github.com/TyphoonMC/go.uuid"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var xpTable = []int{7, 16, 27, 40, 55, 72, 91, 112, 135, 160,
	187, 216, 247, 280, 315, 352, 394, 441, 493, 550, 612,
	679, 751, 828, 910, 997, 1089, 1186, 1288, 1395, 1507,
	1628, 1758, 1897, 2045, 2202, 2368, 2543, 2727, 2920}

type player struct {
	BossBar t.PacketBossBar
	Time int64
	Player *t.Player
	Items []int
	HP float32
	FOOD int
	FOODSAT float32
	XP int
	MONEY int
	BOSS float32
}

func NewPlayer(p *t.Player) *player {
	return &player{
		BossBar: t.PacketBossBar{
			UUID:     uuid.Must(uuid.NewV4()),
			Action:   t.BOSSBAR_ADD,
			Title:    string(config.BossBar),
			Health:   1.0,
			Color:    t.BOSSBAR_COLOR_RED,
			Division: t.BOSSBAR_NODIVISION,
			Flags:    0,
		},
		Time:    time.Now().UnixNano(),
		Player:  p,
		Items:   make([]int, 0),
		HP:      20.0,
		FOOD:    20,
		FOODSAT: 0.0,
		XP:      50,
		MONEY:   100,
		BOSS:    1.0,
	}
}

func (p *player) SetHP(v float32) {
	if v < 0.0 {
		v = 0.0
	}else if v > 20.0 {
		v = 20.0
	}
	p.HP = v
	go p.updateStatus(p.HP, p.FOOD, p.FOODSAT)
}

func (p *player) SetFOOD(v int) {
	if v < 0 {
		v = 0
	}else if v > 20 {
		v = 20
	}
	p.FOOD = v
	go p.updateStatus(p.HP, p.FOOD, p.FOODSAT)
}

func (p *player) SetFOODSAT(v float32) {
	if v < 0.0 {
		v = 0.0
	}else if v > 5.0 {
		v = 5.0
	}
	p.FOODSAT = v
	go p.updateStatus(p.HP, p.FOOD, p.FOODSAT)
}

func (p *player) GetLevel() int {
	var level int
	for j, k := range xpTable {
		level = j
		if p.XP - k < 0 {
			break
		}
	}
	if level == len(xpTable) - 1 {
		level++
	}
	return level
}

func (p *player) SetXP(v int) {
	p.XP = v
	level := p.GetLevel()
	var progress float32
	if level == 0 {
		progress = float32(v) / float32(xpTable[0])
	}else if level == len(xpTable) {
		progress = 0
	}else{
		progress = (float32(v) - float32(xpTable[level-1])) / (float32(xpTable[level]) - float32(xpTable[level-1]))
	}
	go p.updateXP(progress, int64(level), int64(p.XP))
}

func (p *player) Update() {
	go p.updateStatus(p.HP, p.FOOD, p.FOODSAT)
	p.SetXP(p.XP)
}

func (p *player) updateStatus(h float32, f int, fs float32) {
	health := make([]byte, 4)
	binary.BigEndian.PutUint32(health, math.Float32bits(h))
	food := make([]byte, 1)
	binary.PutUvarint(food, uint64(f))
	foodinc := make([]byte, 4)
	binary.BigEndian.PutUint32(foodinc, math.Float32bits(fs))
	pkt := make([]byte, 11)
	pkt = append(pkt, []byte{10, 0x40}...)
	pkt = append(pkt, health...)
	pkt = append(pkt, food...)
	pkt = append(pkt, foodinc...)
	p.Player.WriteRawPacket(pkt)
}

func (p *player) updateXP(b float32, l int64, t int64) {
	bar := make([]byte, 4)
	binary.BigEndian.PutUint32(bar, math.Float32bits(b))
	level := make([]byte, 8)
	x := binary.PutUvarint(level, uint64(l))
	level = level[:x]
	total := make([]byte, 8)
	y := binary.PutUvarint(total, uint64(t))
	total = total[:y]
	pkt := make([]byte, 0)
	pkt = append(pkt, []byte{byte(1 + 4 + x + y), 0x3F}...)
	pkt = append(pkt, bar...)
	pkt = append(pkt, level...)
	pkt = append(pkt, total...)
	p.Player.WriteRawPacket(pkt)
}

func statusHandler(e *t.PlayerChatEvent) {
	uuid := e.Player.GetUUID()
	playersLock.Lock()
	p, ok := players[uuid]
	playersLock.Unlock()
	if p == nil || !ok {
		return
	}
	attack := 0
	shield := 0
	for _, k := range p.Items {
		for _, h := range items {
			if h.Cid != 1 && h.Cid != 2 {
				continue
			}
			if h.Id == k {
				attack += h.Attack
				shield += h.Shield
			}
		}
	}
	output := fmt.Sprintf("Attack: %d ; Shield: %d ; HP: %2.1f ; Hunger: %d ; XP: %d ; MONEY: $%d", attack, shield, p.HP, p.FOOD, p.XP, p.MONEY)
	sendTagMsg(e.Player, "STATUS", output)
}

func itemsHandler(e *t.PlayerChatEvent) {
	uuid := e.Player.GetUUID()
	playersLock.Lock()
	p, ok := players[uuid]
	playersLock.Unlock()
	if p == nil || !ok {
		return
	}
	maps := make(map[string]int, 0)
	sitems := "\n"
	itemlistd := "["
	found := false
	for _, k := range p.Items {
		for _, h := range items {
			if h.Id == k {
				found = true
				itemlistd += fmt.Sprintf("\"[%d] %s %s\", ", h.Id, h.Name, h.EnName)
				key := fmt.Sprintf("[%d] %s %s", h.Id, h.Name, h.EnName)
				if _, ok := maps[key]; !ok {
					maps[key] = 0
				}
				maps[key]++
			}
		}
	}
	if found {
		itemlistd = itemlistd[:len(itemlistd) - 2]
	}
	itemlistd += "]"
	for j, k := range maps {
		sitems += fmt.Sprintf("%s * %d\n", j, k)
	}
	if sitems == "\n" {
		sendTagMsg(e.Player, "ITEMS", "\nempty\nyour full item list:\n" + itemlistd)
	}else{
		sendTagMsg(e.Player, "ITEMS", sitems + "\nyour full item list:\n" + itemlistd)
	}
}

func useHandler(e *t.PlayerChatEvent) {
	uuid := e.Player.GetUUID()
	playersLock.Lock()
	p, ok := players[uuid]
	playersLock.Unlock()
	if p == nil || !ok {
		return
	}
	sep := strings.Split(e.Message, " ")
	if len(sep) != 2 {
		sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
		return
	}
	if sep[0] != "use" {
		sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
		return
	}
	n, err := strconv.Atoi(sep[1])
	if err != nil {
		sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
		return
	}
	match := false
	for _, k := range items {
		if k.Cid != 3 && k.Cid != 4 {
			continue
		}
		if k.Id == n {
			match = true
			break
		}
	}
	if !match {
		sendTagMsg(e.Player, "SYSTEM", "\ninvalid item")
		return
	}
	if n == 15 {
		match := false
		for _, k := range p.Items {
			if k == 6 {
				match = true
				break
			}
		}
		if !match {
			sendTagMsg(e.Player, "SYSTEM", "\nCan't use arrow without bow.")
			return
		}
	}
	match = false
	for j, k := range p.Items {
		if k == n {
			match = true
			p.Items, _ = slicePop(p.Items, uint(j))
			break
		}
	}
	if !match {
		sendTagMsg(e.Player, "SYSTEM", "\nyou must buy the item first")
		return
	}
	attack := 0
	hp := float32(0)
	food := 0
	xp := 0
	for _, h := range items {
		if h.Cid != 3 && h.Cid != 4 {
			continue
		}
		if n == 15 && h.Id == 6 {
			attack += h.Attack
			hp += h.HP
			food += h.Food
			xp += h.XP
		}
		if h.Id == n {
			attack += h.Attack
			hp += h.HP
			food += h.Food
			xp += h.XP
		}
	}
	sattack := float32(attack) / 10 / 100
	playersLock.Lock()
	if players[uuid].FOOD + food < 0 {
		playersLock.Unlock()
		sendTagMsg(e.Player, "SYSTEM", "\nuse fail. low Hunger.")
		return
	}
	players[uuid].HP += hp
	if players[uuid].HP < 0 {
		players[uuid].HP = 0
	}
	myhp := players[uuid].HP
	players[uuid].FOOD += food
	if players[uuid].FOOD < 0 {
		players[uuid].FOOD = 0
	}
	players[uuid].XP += xp
	if players[uuid].XP < 0 {
		players[uuid].XP = 0
	}
	bosshp := players[uuid].BOSS
	if attack > 0 {
		players[uuid].BOSS -= sattack
		if players[uuid].BOSS < 0 {
			players[uuid].BOSS = 0
		}
		bosshp = players[uuid].BOSS
		go setBossHealth(players[uuid].Player, players[uuid].BOSS)
	}
	go players[uuid].Update()
	playersLock.Unlock()
	if attack <= 0 {
		sendTagMsg(e.Player, "SYSTEM", fmt.Sprintf("\nuse succ."))
		go statusHandler(e)
	}else{
		sendTagMsg(e.Player, "SYSTEM", fmt.Sprintf("\nuse succ.\nBOSS's HP -%d. [%d / 1000]", attack, int(bosshp * 1000)))
		go statusHandler(e)
	}
	if myhp > 0 && bosshp == 0 {
		defeatBoss(e)
	}
}

func attackHandler(e *t.PlayerChatEvent) {
	uuid := e.Player.GetUUID()
	playersLock.Lock()
	p, ok := players[uuid]
	playersLock.Unlock()
	if p == nil || !ok {
		return
	}
	attack := 0
	shield := 0
	hp := float32(0)
	food := 0
	xp := 0
	match := false
	for _, k := range p.Items {
		for _, h := range items {
			if h.Cid != 1 && h.Cid != 2 {
				continue
			}
			if h.Id == 6 {
				continue
			}
			if h.Id == k {
				attack += h.Attack
				shield += h.Shield
				hp += h.HP
				food += h.Food
				xp += h.XP
				if h.Cid == 1 {
					match = true
				}
			}
		}
	}
	if attack == 0 {
		attack = 1
	}
	sattack := float32(attack) / 10 / 100
	beattack := float32(rand.Intn(500)) / 100 + 0.5 - float32(shield) / 100
	playersLock.Lock()
	if players[uuid].HP + hp < 0 {
		playersLock.Unlock()
		sendTagMsg(e.Player, "SYSTEM", "\nattack fail. low HP.")
		return
	}
	if players[uuid].FOOD + food < 0 {
		playersLock.Unlock()
		sendTagMsg(e.Player, "SYSTEM", "\nattack fail. low Hunger.")
		return
	}
	players[uuid].HP += hp
	if beattack > 0 {
		players[uuid].HP -= beattack
	}
	if players[uuid].HP < 0 {
		players[uuid].HP = 0
	}
	myhp := players[uuid].HP
	players[uuid].FOOD += food
	if players[uuid].FOOD < 0 {
		players[uuid].FOOD = 0
	}
	players[uuid].XP += xp
	if players[uuid].XP < 0 {
		players[uuid].XP = 0
	}
	players[uuid].BOSS -= sattack
	if players[uuid].BOSS < 0 {
		players[uuid].BOSS = 0
	}
	bosshp := players[uuid].BOSS
	go setBossHealth(players[uuid].Player, players[uuid].BOSS)
	go players[uuid].Update()
	playersLock.Unlock()
	if !match {
		sendTagMsg(e.Player, "SYSTEM", "\nwithout weapon, attacking...")
		go statusHandler(e)
	}else{
		sendTagMsg(e.Player, "SYSTEM", fmt.Sprintf("\nattack succ.\nBOSS's HP -%d. [%d / 1000]", attack, int(bosshp * 1000)))
		go statusHandler(e)
	}
	if myhp > 0 && bosshp == 0 {
		defeatBoss(e)
	}
}

func defeatBoss(e *t.PlayerChatEvent) {
	sendTagMsg(e.Player, "ADMIN", "\nCongratulation!\nEncoded Message:\nF5GUGMRQ\nGIYC2RCF\nIJKUOLKW\nJFCVOORN\nFEFFC4RR\nKBDVG62G\nGNYGQZJT\nL5EGMTSU\nGNPTA4ZN\nKRBF6RTZ\nOZYHE7T5\n")
	core.GetPlayerRegistry().ForEachPlayerAsync(func(player *t.Player){
		msg := t.ChatMessage("")
		msg.SetExtra([]t.IChatComponent{
			t.ChatMessage("Congratulation! <"),
			t.ChatMessage(e.Player.GetName()),
			t.ChatMessage("> defeat the BOSS!"),
		})
		player.SendActionBar(msg)
	})
}

func debugViewHandler(e *t.PlayerChatEvent) {
	if e.Message == "MC2020-DEBUG-VIEW:-)" {
		sendTagMsg(e.Player, "ADMIN", "\nUsage: /MC2020-DEBUG-VIEW:-) [Player_UUID] [File_Offset_In_Bytes]\n      Read player's log file. Output is base64 encoded. Maximum reading of raw bytes is 1 MB (1048576 Bytes).")
		return
	}else if e.Message[:20] == "MC2020-DEBUG-VIEW:-)" {
		sep := strings.Split(e.Message, " ")
		if len(sep) != 3 {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		n, err := strconv.ParseInt(sep[2], 10, 64)
		if err != nil || n < 0 {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		if sep[0] != "MC2020-DEBUG-VIEW:-)" {
			sendTagMsg(e.Player, "SYSTEM", "\ninvalid params")
			return
		}
		f, err := os.OpenFile(fmt.Sprintf("logs/%s", sep[1]), os.O_RDONLY, 0400)
		if err != nil {
			sendTagMsg(e.Player, "ADMIN", "\nlog file not found")
			return
		}
		n, err = f.Seek(n, 0)
		if err != nil {
			sendTagMsg(e.Player, "SYSTEM", "\ninternal error")
			return
		}
		buf := make([]byte, 1048576)
		x, err := f.Read(buf)
		f.Close()
		if err != nil {
			sendTagMsg(e.Player, "DEBUG-VIEW", "\nEOF")
			return
		}
		buf = buf[:x]
		str := base64.StdEncoding.EncodeToString(buf)
		sendTagMsg(e.Player, "DEBUG-VIEW", "\n" + str)
		return
	}
}