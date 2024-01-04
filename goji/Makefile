setup-db:
	@sudo docker run -it -p 5432:5432 --name goji -e POSTGRES_USER=cube -e POSTGRES_PASSWORD=sec postgres:13

inspect-db:
	@sudo docker inspect goji

start-db:
	@sudo docker start goji

stop-db:
	@sudo docker stop goji
