from golang:1.11
RUN mkdir -p /yt-downloader
WORKDIR /yt-downloader
RUN apt-get update && apt-get install -y python3 python3-pip lsof
ADD . .
RUN pip3 install -r python/requirements.txt
WORKDIR go 
RUN go get -d && go build -o run .
WORKDIR /yt-downloader
VOLUME /Music /url-list.txt
ENTRYPOINT bash start.sh
