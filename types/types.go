package types

import (
	"encoding/json"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

type Auction struct {
	Apr           float64
	EndTime       int64
	HighestBidder int64
}

type SocketReq struct {
	Apr    float64 `json:"apr"`
	Uuid   int64   `json:"uuid"`
	Bidder int64   `json:"bidder"`
}

type AuctionCreateReq struct {
	BondId   int64   `json:"bondId"`
	SellerId int64   `json:"sellerId"`
	MaxApr   float64 `json:"maxApr"`
}

type AuctionCreateRes struct {
	Valid bool `json:"valid"`
}

type AuctionEndReq struct {
	Auction  Auction `json:"auction"`
	BondUuid int64   `json:"bondUuid"`
}

func (a *Auction) Bid(req SocketReq) (bool, bool) {
	if req.Apr < a.Apr && time.Now().UnixNano() > a.EndTime {
		a.Apr = req.Apr
		a.HighestBidder = req.Bidder

		return true, false
	} else if time.Now().UnixNano() > a.EndTime {
		body, err := json.Marshal(AuctionEndReq{
			Auction:  *a,
			BondUuid: req.Uuid,
		})

		if err != nil {
			return false, false
		}

		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()
		req.AppendBody(body)
		req.SetRequestURI(os.Getenv("DOMAIN_NAME") + "/api/finishAuction")

		err = Client.Do(req, res)

		if err != nil {
			return false, false
		}

		return false, true
	}

	return false, false
}
