RUN Python
docker inspect <container_name> | python docker-compose-converter.py - > docker-compose.yml
or
docker inspect <container_name> | python3 docker-compose-converter.py - > docker-compose.yml

RUN GO
./docker-compose-converter --container my-container --output docker-compose.yml
