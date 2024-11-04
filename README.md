<h4>RUN Python</h4><br>
    <h5>docker-compose-converter</h5> 
       docker inspect <container_name> | python docker-compose-converter.py - > docker-compose.yml`
   
<h4>RUN GO</h4><br>
    <h5>docker-compose-converter</h5>
       ./docker-compose-converter --container my-container --output docker-compose.yml

    <h5>cryptdecrypt</h5><br>
       <code>./cryptdecrypt -mode crypt -password -text</code>
       ./cryptdecrypt -mode decrypt -password -text salt:ciphertext
