RUN Python<br>
<code>
docker inspect <container_name> | python docker-compose-converter.py - > docker-compose.yml
</code>
or
<code>
docker inspect <container_name> | python3 docker-compose-converter.py - > docker-compose.yml
</code>
<br><br>
RUN GO<br>
<code>
./docker-compose-converter --container my-container --output docker-compose.yml
</code>
