setup-db:
	@sudo docker run -it -p 5432:5432 --name goji -e POSTGRES_USER=cube -e POSTGRES_PASSWORD=sec postgres:13

inspect-db:
	@sudo docker inspect goji

inspect-with-curl:
	@docker ps | grep -e "goji" | awk '{ print $1 }' | xargs -I {} curl --unix-socket /var/run/docker.sock http://docker/containers/\{\}/json | jq .

start-db:
	@sudo docker start goji

stop-db:
	@sudo docker stop goji
