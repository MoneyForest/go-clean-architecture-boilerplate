.PHONY: docker-up
docker-up:
	docker compose up -d --build

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: docker-clean
docker-clean:
	docker compose down -v
	sudo rm -rf ./docker/mysql/data/*
	sudo mkdir -p ./docker/mysql/data
	sudo touch ./docker/mysql/data/.keep
	sudo chmod 777 ./docker/mysql/data
	sudo chown -R 999:999 ./docker/mysql/data
	find ./docker/redis/data -type f ! -name '.keep' -delete
	find ./docker/localstack/data -type f ! -name '.keep' -delete

.PHONY: docker-down-and-up
docker-down-and-up: docker-down docker-up

.PHONY: docker-clean-and-up
docker-clean-and-up: docker-clean docker-up

.PHONY: docker-logs
docker-logs:
	docker compose logs -f

.PHONY: docker-setup
docker-setup: docker-clean-and-up
