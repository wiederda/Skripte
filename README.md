RUN Python<br>
docker inspect <container_name> | python docker-compose-converter.py - > docker-compose.yml
br>
or
<br>>
docker inspect <container_name> | python3 docker-compose-converter.py - > docker-compose.yml
<br><br>
RUN GO<br>
./docker-compose-converter --container my-container --output docker-compose.yml
