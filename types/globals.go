package types

import "github.com/fasthttp/websocket"

var Sockets map[string][]*websocket.Conn
var BidMap map[string]Auction

func Init() {
	Sockets = make(map[string][]*websocket.Conn)
	BidMap = make(map[string]Auction)
}
