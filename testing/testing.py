import requests
import asyncio
from json import dumps
from websockets.sync.client import connect

def sendBid(Apr, Uuid, Bidder):
    with connect("ws://localhost:8001/connect") as ws:
        data = dumps(
            dict(
                apr = Apr,
                uuid = Uuid,
                bidder = Bidder
            )
        )
        ws.send(data)
        msg = ws.recv()
        print(msg)

def createDummyAuction(BondId, SellerId, MaxApr):
    args = dict(
        bondId = BondId,
        sellerId = SellerId,
        maxApr = MaxApr
    )
    data = requests.post("http://localhost:8001/create_dummy_auction", data = dumps(args))
    print(data.content)


if __name__ == "__main__":
    createDummyAuction(0, 0, 10000)
    sendBid(1, 0, 2)

    


