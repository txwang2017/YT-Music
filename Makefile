build:	
	docker build -t yt-downloader .

all: build

.PHONY: build all
