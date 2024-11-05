#### RUN Python

#### docker-compose-converter

`docker inspect <container_name> | python docker-compose-converter.py - > docker-compose.yml`

#### RUN GO

##### docker-compose-converter

##### -container string
   #####     Name or ID of the Docker container (optional, only used if input file is not provided)
##### -input string
   #####     Path to the input YAML file with container information
##### -output string
   #####     Path to the output Docker Compose file (default "docker-compose.yml")

` ./docker-compose-converter -container my-container -output docker-compose.yml `


` ./docker-compose-converter -input -output `

##### cryptdecrypt
 `./cryptdecrypt -mode crypt -password -text ` <br>
` ./cryptdecrypt -mode decrypt -password -text salt:ciphertext `
