**RUN Python**  
    *docker-compose-converter*    
       docker inspect <container_name> | python docker-compose-converter.py - > docker-compose.yml`
   
**RUN GO**  
    *docker-compose-converter*
       ./docker-compose-converter --container my-container --output docker-compose.yml

    *cryptdecrypt*
       ./cryptdecrypt crypt oder decrypt -password -text   
