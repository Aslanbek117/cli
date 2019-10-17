# cli
cli written in Go

#How to test
- clone repo
- cd to folder with Dockerfile
- docker file build -t test . 
- docker run -i -t test /bin/bash
- in container: app -url="someUrl" -pattern=world

