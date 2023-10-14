import requests
import asyncio
from json import dumps
from websockets.sync.client import connect

def sendBid():
    r = requests.get("ws://localhost:8001/connect")
    if (r.status_code == 200) :
        with connect("ws://localhost:8001/connect") as ws:
            ws.send("Hello world!")

def createDummyAuction(BondId, SellerId, MaxApr):
    args = dict(
        bondId = BondId,
        sellerId = SellerId,
        maxApr = MaxApr
    )
    data = requests.post("http://localhost:8001/create_dummy_auction", data = dumps(args))
    print(data.content)


if __name__ == "__main__":
    createDummyAuction(0, 0, 0.5)

    


