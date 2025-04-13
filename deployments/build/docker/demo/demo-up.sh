#1/bin/sh

mkdir -p /opt/postgres
docker compose -f docker-compose.yaml pull
docker compose -f docker-compose.yaml up
docker compose -f docker-compose.yaml down

