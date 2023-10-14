package types

var Sockets map[int64][]SocketClient
var BidMap map[int64]Auction

func Init() {
	Sockets = make(map[int64][]SocketClient)
	BidMap = make(map[int64]Auction)
}
