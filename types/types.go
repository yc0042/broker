package types

import "time"

type Auction struct {
	Apr           float32
	endTime       int64
	HighestBidder int64
}

type SocketReq struct {
	Apr    float32 `json:apr`
	Uuid   string  `json:uuid`
	Bidder int64   `json:bidder`
}

func (a *Auction) Bid(req SocketReq) bool {
	if req.Apr < a.Apr && time.Now().UnixNano() > a.endTime {
		a.Apr = req.Apr
		a.HighestBidder = req.Bidder

		return true
	}

	return false
}
