import socket
#For test only
s = socket.socket()
s.connect(('127.0.0.1', 8000))
s.sendall('{"id": 19, "url": "https://www.youtube.com/watch?v=VC2rAxRID9s"}'.encode('utf-8'))
data = []
while True:
    buff = s.recv(1024)
    if not buff:
        break
    data.extend(buff)
data = ''.join(map(lambda x: chr(x), data))
