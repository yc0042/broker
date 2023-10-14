package types

import "github.com/fasthttp/websocket"

var Sockets map[int64][]*websocket.Conn
var BidMap map[int64]Auction

func Init() {
	Sockets = make(map[int64][]*websocket.Conn)
	BidMap = make(map[int64]Auction)
}
