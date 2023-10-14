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
    data = requests.post(url = "https://localhost:8001/create_dummy_auction", params = args)
    print(data.content)

def 

    


