package main

import (
	"lendshare/broker/types"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/connect":
		err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
			var req types.SocketReq
			first := true
			for {
				err := conn.ReadJSON(req)
				if err != nil {
					conn.Close()
					break
				}
				if first {
					types.Sockets[req.Uuid] = append(types.Sockets[req.Uuid], conn)
					first = false
				}
				auction := types.BidMap[req.Uuid]
				s := auction.Bid(req)

				if s {
					for _, socket := range types.Sockets[req.Uuid] {
						socket.WriteJSON(types.SocketReq{
							Apr:    auction.Apr,
							Bidder: auction.HighestBidder,
						})
					}
				}
			}
		})

		if err != nil {
			panic(err)
		}
	}
}

func main() {
	types.Init()
	fasthttp.ListenAndServe("8001", handler)
}
