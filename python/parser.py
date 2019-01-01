import asyncio, json
from urllib.parse import unquote

import aiohttp
from bs4 import BeautifulSoup

class Parser:
    def __init__(self, data):
        obj = json.loads(data.decode('utf-8'))
        self.id = obj.get('id')
        self.url = obj.get('url')
        self.headers = {
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
            "accept-language": "en-US,en;q=0.9",
            "cookie": "VISITOR_INFO1_LIVE=7gcjyw3ehOI; CONSENT=YES+US.en+20170326-06-0; PREF=f1=50000000; GPS=1; YSC=lTe68erqg-M",
            "upgrade-insecure-requests": "1",
            "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
        }
        self.request = aiohttp.request("GET", self.url, headers=self.headers)
    
    def _get_raw_data(self, raw: str) -> [str]:
        ret = []
        count = 0
        tail = 0
        head = 0
        while head < len(raw):
            if raw[head] == '{':
                count += 1
                tail = head
                while tail + 1 < len(raw) and count > 0:
                    tail += 1
                    if raw[tail] == '}':
                        count -= 1
                    elif raw[tail] == '{':
                        count += 1
                ret.append(raw[head:tail+1])
                head = tail + 1
            else:
                head += 1
        return ret
    
    def _parse(self):
        soup = BeautifulSoup(self.content, "html.parser")
        scripts = soup.find(name="div", attrs={"id": "player-wrap"}).find_all(name="script")
        script = scripts[1].get_text()
        raw_data = self._get_raw_data(script)[1]
        raw_data = json.loads(raw_data).get('args', {}).get('player_response', "{}")
        self.data = json.loads(raw_data).get('streamingData', {}).get('adaptiveFormats', [])
        for d in self.data:
            if 'audio/mp4' in d.get('mimeType'):
                d['mimeType'] = unquote(d['mimeType'])
                return d
        return None

    async def get_audio_info(self):
        async with self.request as resp:
            self.content = await resp.read()
        audio_info = self._parse()
        audio_info.update({'id': self.id})
        return json.dumps(audio_info)
