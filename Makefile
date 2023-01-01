create_rabbit:
	docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-alpine
run_rabbit:
	docker start some-rabbit
stop_rabbit:
	docker stop some-rabbit
run_producer:
	go run cmd/producer/main.go
run_consumer_ping:
	go run cmd/consumer/main.go Ping
run_consumer_pong:
	go run cmd/consumer/main.go Pong

.PHONY: create_rabbit run_rabbit stop_rabbit run_producer run_consumer_ping run_consumer_pong