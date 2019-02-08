python3 python/server.py 2>&1 1>/dev/null &
while true; do
	ready=$(lsof -i:8000)
	if [[ ! -z $ready ]]; then 
		break 
	fi
	sleep 1
done;
./go/run -url "$url" -name "$name" -list "$list"
