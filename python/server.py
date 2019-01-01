# import socket, asyncio

# port = 8000


# async def server(reader, writer):
#     data = await reader.read(1024)



# if __name__ == '__main__':
#     loop = asyncio.new_event_loop()
#     coro = asyncio.start_server(server, '127.0.0.1', 8000)
#     loop.run_until_complete(coro)
#     loop.run_forever()