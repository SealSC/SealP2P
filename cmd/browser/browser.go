package main

import (
	"github.com/SealSC/SealP2P"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/SealSC/SealP2P/conn/msg"
	"github.com/gorilla/websocket"
	"log"
	_ "embed"
	"github.com/SealSC/SealP2P/conf"
)

type wsMSG struct {
	ws *websocket.Conn
}

func (w *wsMSG) OnMessage(p *msg.Message) *msg.Message {
	err := w.ws.WriteJSON(p)
	if err != nil {
		panic(err)
	}
	return nil
}

//go:embed index.html
var indexHTML []byte

func main() {
	err := SealP2P.InitLocalNode(conf.DefaultConfig)
	if err != nil {
		panic(err)
	}
	node := SealP2P.LocalNode()
	engine := gin.New()
	log.Println("node id:", node.GetNodeID())
	engine.GET("/", func(c *gin.Context) {
		c.Writer.Write(indexHTML)
	})

	engine.Any("/log", func(c *gin.Context) {
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		ws.WriteMessage(websocket.TextMessage, []byte("show..........."))
		node.SetMessenger(&wsMSG{ws: ws})
		select {}
	})

	engine.Any("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, node.GetNodeStatus())
	})

	engine.Any("/join", func(c *gin.Context) {
		if err := node.Join(); err != nil {
			panic(err)
		}
		c.Status(http.StatusOK)
	})
	engine.Any("/leave", func(c *gin.Context) {
		if err := node.Leave(); err != nil {
			panic(err)
		}
		c.Status(http.StatusOK)
	})
	engine.Any("/multicast", func(c *gin.Context) {
		node.MulticastMsg(readForm(c))
		c.Status(http.StatusOK)
	})
	engine.Any("/broadcast", func(c *gin.Context) {
		if err := node.BroadcastMsg(readForm(c)); err != nil {
			panic(err)
		}
		c.Status(http.StatusOK)
	})
	engine.Any("/send", func(c *gin.Context) {
		if err := node.SendMsg(readForm(c)); err != nil {
			panic(err)
		}
		c.Status(http.StatusOK)
	})
	engine.Any("/nodes", func(c *gin.Context) {
		list := node.GetNodeList()
		c.JSON(http.StatusOK, list)
	})

	if err := engine.Run(":8080"); err != nil {
		panic(err)
	}
}

func readForm(c *gin.Context) *msg.Message {
	body := c.PostForm("body")
	tos := strings.Split(c.PostForm("tos"), "\n")
	payload := SealP2P.EmptyMessage()
	payload.Payload = []byte(body)
	payload.ToID = tos
	return payload
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
