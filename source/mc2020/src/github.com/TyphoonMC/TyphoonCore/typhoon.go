package typhoon

import (
	"bufio"
	"github.com/TyphoonMC/go.uuid"
	"log"
	"math/rand"
	"net"
	"reflect"
	"time"
)

type Core struct {
	connCounter      int
	eventHandlers    map[reflect.Type][]EventCallback
	brand            string
	rootCommand      CommandNode
	compiledCommands []commandNode
	playerRegistry   *PlayerRegistry
}

func Init() *Core {
	initConfig()
	initPackets()
	initHacks()
	c := &Core{
		0,
		make(map[reflect.Type][]EventCallback),
		"MC2020",
		CommandNode{
			commandNodeTypeRoot,
			nil,
			nil,
			nil,
			"",
			nil,
		},
		nil,
		newPlayerRegistry(),
	}
	c.compileCommands()
	return c
}

func (c *Core) Start() {
	ln, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server launched on port", config.ListenAddress)
	go c.keepAlive()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		} else {
			c.connCounter += 1
			go c.handleConnection(conn, c.connCounter)
		}
	}
}

func (c *Core) SetBrand(brand string) {
	br := make([]byte, len(brand)+1)
	copy(br[:len(brand)], []byte(brand))
	c.brand = string(br)
}

func (c *Core) GetPlayerRegistry() *PlayerRegistry {
	return c.playerRegistry
}

func (c *Core) keepAlive() {
	r := rand.New(rand.NewSource(15768735131534))
	keepalive := &PacketPlayKeepAlive{
		Identifier: 0,
	}
	for {
		c.playerRegistry.ForEachPlayer(func(player *Player) {
			if player.state == PLAY {
				if player.keepalive != 0 {
					return
				}

				id := int(r.Int31())
				keepalive.Identifier = id
				player.keepalive = id
				player.WritePacket(keepalive)
				go func() {
					nowTime := time.Now().Unix()
					for {
						if player.keepalive == 0 {
							break
						}
						if time.Now().After(time.Unix(nowTime, 0).Add(40 * time.Second)) {
							player.Kick("Timed out")
							break
						}
						time.Sleep(1 * time.Second)
					}
				}()
			}
		})
		time.Sleep(5 * time.Second)
	}
}

func (c *Core) handleConnection(conn net.Conn, id int) {
	log.Printf("%s(#%d) connected.", conn.RemoteAddr().String(), id)

	_uuid, err := uuid.NewV4()
	if err != nil {
		conn.Close()
		return
	}

	player := &Player{
		core:     c,
		id:       id,
		conn:     conn,
		state:    HANDSHAKING,
		protocol: V1_10,
		io: &ConnReadWrite{
			rdr: bufio.NewReader(conn),
			wtr: bufio.NewWriter(conn),
		},
		inaddr: InAddr{
			"",
			0,
		},
		name:        "",
		uuid:        _uuid.String(),
		keepalive:   0,
		compression: false,
	}

	for {
		_, err := player.ReadPacket()
		if err != nil {
			break
		}
	}

	if player.state == PLAY {
		player.core.CallEvent(&PlayerQuitEvent{player})
		player.unregister()
	}
	conn.Close()
	log.Printf("%s(#%d) disconnected.", conn.RemoteAddr().String(), id)
}
