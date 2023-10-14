package types

import "time"

type Auction struct {
	Apr     float32
	endTime int64
}

func (a *Auction) Bid(apr float32) {
	if apr < a.Apr && time.Now().UnixNano() > a.endTime {
		a.Apr = apr
	}
}
