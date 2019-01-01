import socket, asyncio

from parser import Parser

port = 8000


async def server(reader, writer):
    data = await reader.read(1024)
    parser = Parser(data)
    audio_info = await parser.get_audio_info()
    writer.write(audio_info.encode('utf-8'))
    await writer.drain()
    writer.close()

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    coro = asyncio.start_server(server, '127.0.0.1', 8000)
    loop.run_until_complete(coro)
    loop.run_forever()