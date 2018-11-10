build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shjp-dao .

docker_update: build
	docker build --no-cache --tag dao_server .
	docker run --name shjp_dao_server -itd dao_server
	docker commit shjp_dao_server sanggonlee/shjp:shjp_dao_server
	docker stop shjp_dao_server
	docker rm shjp_dao_server
	docker push sanggonlee/shjp:shjp_dao_server
