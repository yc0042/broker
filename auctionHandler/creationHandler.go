package auctionhandler

import (
	"encoding/json"
	"fmt"
	"lendshare/broker/types"
	"os"

	"github.com/valyala/fasthttp"
)

func CreateAuction(ctx *fasthttp.RequestCtx) (types.AuctionCreateReq, error) {
	var reqBody types.AuctionCreateReq

	err := json.Unmarshal(ctx.Request.Body(), &reqBody)

	if err != nil {
		return types.AuctionCreateReq{}, err
	}

	req := fasthttp.AcquireRequest()
	req.AppendBody(ctx.Request.Body())
	req.SetRequestURI(os.Getenv("DOMAIN_NAME") + "/api/create_auction")
	res := fasthttp.AcquireResponse()
	err = types.Client.Do(req, res)

	if err != nil {
		return types.AuctionCreateReq{}, err
	}
	var resBody types.AuctionCreateRes
	err = json.Unmarshal(res.Body(), &resBody)

	if err != nil {
		return types.AuctionCreateReq{}, err
	}

	if resBody.Valid {
		return reqBody, nil
	}

	return types.AuctionCreateReq{}, fmt.Errorf("auction invalid")
}
