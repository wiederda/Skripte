#### RUN Python

#### docker-compose-converter

`docker inspect <container_name> | python docker-compose-converter.py - > docker-compose.yml`

#### RUN GO

##### docker-compose-converter

` ./docker-compose-converter -container my-container -output docker-compose.yml `
optional, on Docker Host, only used if input file is not provided

` ./docker-compose-converter -input -output `

##### cryptdecrypt
 `./cryptdecrypt -mode crypt -password -text ` <br>
` ./cryptdecrypt -mode decrypt -password -text salt:ciphertext `
