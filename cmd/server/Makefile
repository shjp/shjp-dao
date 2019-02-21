build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shjp-dao .

docker_update: build
	docker build --no-cache --tag shjp-dao-server_v$(version) .
	docker run --name shjp-dao-server -itd shjp-dao-server_v$(version)
	docker commit shjp-dao-server sanggonlee/shjp:shjp-dao-server_v$(version)
	docker stop shjp-dao-server
	docker rm shjp-dao-server
	docker push sanggonlee/shjp:shjp-dao-server_v$(version)
	rm ./shjp-dao

docker_update_dev: build
	docker build --no-cache --tag shjp-dao-server_dev .
	docker run --name shjp-dao-server -itd shjp-dao-server_dev
	docker commit shjp-dao-server sanggonlee/shjp:shjp-dao-server_dev
	docker stop shjp-dao-server
	docker rm shjp-dao-server
	docker push sanggonlee/shjp:shjp-dao-server_dev
	rm ./shjp-dao