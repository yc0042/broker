package main

import (
	"encoding/json"
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
					delete(types.BidMap, req.Uuid)
				}
			}
		})

		if err != nil {
			panic(err)
		}

	case "/create_dummy_auction":
		var body types.AuctionCreateReq
		json.Unmarshal(ctx.Request.Body(), &body)
		types.BidMap[body.BondId] = types.Auction{
			Apr:     body.MaxApr,
			EndTime: time.Now().AddDate(0, 0, 1).UnixNano(),
		}

		ctx.SetStatusCode(200)
		ctx.Response.AppendBodyString("Success creating auction")
	case "/create_auction":
		info, err := auctionhandler.CreateAuction(ctx)

		if err != nil {
			ctx.SetStatusCode(400)
			ctx.Response.AppendBodyString("Failure creating auction")
		} else {
			tomorrow := time.Now().AddDate(0, 0, 1).UnixNano()
			types.BidMap[info.BondId] = types.Auction{
				Apr:     info.MaxApr,
				EndTime: tomorrow,
			}
			ctx.SetStatusCode(200)
			ctx.Response.AppendBodyString("Success creating auction")
		}

	case "/get_auctions":
		var body types.BatchAuctionReq
		err := json.Unmarshal(ctx.Request.Body(), &body)

		var out []types.AuctionEndReq

		out = make([]types.AuctionEndReq, 0)

		if err != nil {
			ctx.SetStatusCode(400)
		} else {
			if len(body.Uuids) == 0 {
				for uuid, auction := range types.BidMap {
					out = append(out, types.AuctionEndReq{
						BondUuid: uuid,
						Auction:  auction,
					})
				}
			} else {
				for _, uuid := range body.Uuids {
					out = append(out, types.AuctionEndReq{
						BondUuid: uuid,
						Auction:  types.BidMap[uuid],
					})
				}
			}
		}

		res := types.BatchAuctionRes{
			Auctions: out,
		}

		ser, err := json.Marshal(res)

		if err != nil {
			ctx.SetStatusCode(500)
		} else {
			ctx.Response.AppendBody(ser)
		}
	}
}

func main() {
	godotenv.Load()
	types.InitClient()
	types.Init()
	fasthttp.ListenAndServe("8001", handler)
}
