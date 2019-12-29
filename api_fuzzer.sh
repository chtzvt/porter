PORTER_HOST='https://garage.local:8080'
API_KEY='changeme'

API_ENDPOINTS=( "/api/v1/list"
	   "/api/v1/state/all"
	   "/api/v1/state/test"
           "/api/v1/lock/test"
           "/api/v1/open/test"
           "/api/v1/close/test"
           "/api/v1/trip/test"
           "/api/v1/unlock/test")

CMD_NUM=1
while true;
do
	CMD=$PORTER_HOST`shuf -n1 -e "${API_ENDPOINTS[@]}"`
	echo ">>> Fuzzer Command #$CMD_NUM -- $CMD"
	curl -H "Authorization: Bearer $API_KEY" -k $CMD
	CMD_NUM=$((CMD_NUM+1))
done
