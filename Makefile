include .env
export

service-run:
	@go run main.go


service-deploy:
	docker compose up -d

service-undeploy:
	docker compose down
