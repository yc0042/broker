package main

import (
	auctionhandler "lendshare/broker/auctionHandler"
	"lendshare/broker/types"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/joho/godotenv"
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
				validBid, auctionEnded := auction.Bid(req)

				if validBid {
					for _, socket := range types.Sockets[req.Uuid] {
						socket.WriteJSON(types.SocketReq{
							Apr:    auction.Apr,
							Bidder: auction.HighestBidder,
						})
					}
				}

				if auctionEnded {
					for _, socket := range types.Sockets[req.Uuid] {
						socket.Close()
					}

					delete(types.Sockets, req.Uuid)
				}
			}
		})

		if err != nil {
			panic(err)
		}

	case "/create_auction":
		info, err := auctionhandler.CreateAuction(ctx)

		if err != nil {
			ctx.SetStatusCode(400)
		}

		tomorrow := time.Now().AddDate(0, 0, 1).UnixNano()
		types.BidMap[info.BondId] = types.Auction{
			Apr:     info.MaxApr,
			EndTime: tomorrow,
		}
	}
}

func main() {
	godotenv.Load()
	types.InitClient()
	types.Init()
	fasthttp.ListenAndServe("8001", handler)
}
