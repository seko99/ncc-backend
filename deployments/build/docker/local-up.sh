#1/bin/sh

mkdir -p /opt/postgres
mkdir -p /opt/rabbitmq
docker compose -f docker-compose.yaml pull
docker compose -f docker-compose.yaml up
docker compose -f docker-compose.yaml down

