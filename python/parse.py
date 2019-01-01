import asyncio, json

import aiohttp

class Parse:
    def __init__(self, data):
        obj = json.loads(data)
        self.id = obj.get('id')
        self.url = obj.get('url')
        self.headers = {
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
            "accept-language": "en-US,en;q=0.9",
            "cookie": "VISITOR_INFO1_LIVE=7gcjyw3ehOI; CONSENT=YES+US.en+20170326-06-0; PREF=f1=50000000; GPS=1; YSC=lTe68erqg-M",
            "upgrade-insecure-requests": "1",
            "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
        }
        self.request = aiohttp.request("GET", self.url, headers=json.dumps(self.headers))
    
    async def request(self):
        async with self.request as rep:
            r = await rep.
            print(r)


if __name__ == "__main__":
    data = '{"id": 10, "url": "http://www.google.com"}'
    p = Parse(data)
    loop = asyncio.get_event_loop()
    loop.run_until_complete(p.request())