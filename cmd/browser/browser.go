package main

import (
	"github.com/SealSC/SealP2P"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/SealSC/SealP2P/conn/msg"
	"log"
)

func main() {
	node := SealP2P.LocalNode()
	engine := gin.New()
	log.Println("node id:", node.GetNodeID())
	engine.StaticFile("/", "index.html")
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

func readForm(c *gin.Context) *msg.Payload {
	body := c.PostForm("body")
	tos := strings.Split(c.PostForm("tos"), "\n")
	payload := SealP2P.NewPayload()
	payload.Body = []byte(body)
	payload.ToID = tos
	return payload
}
