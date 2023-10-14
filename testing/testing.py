import requests
import asyncio
from json import dumps
from websockets.sync.client import connect

def sendBid():
    r = requests.get("ws://localhost:8001/connect")
    if (r.status_code == 200) :
        with connect("ws://localhost:8001/connect") as ws:
            ws.send("Hello world!")



    


