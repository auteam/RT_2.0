#!/bin/bash
GREEN='\033[0;32m'
NC='\033[0m'
#Switching dev and prod branches
if [[ $1 == "switch" && $2 == "dev" ]] 
then
	printf "Switching docker (Frontend) to ${GREEN}DEV${NC} branch...\n"
	cp Dockerfiles/Dockerfile.frontend.dev frontend/Dockerfile
	cp Dockerfiles/nginx.conf.dev nginx/nginx.conf
	cp Dockerfiles/docker-compose.yml.dev ./docker-compose.yml
	printf "${GREEN}Done!${NC}\n"
fi
if [[ $1 == "switch" && $2 == "prod" ]]
then
	printf "Switching docker (Frontend) to ${GREEN}PROD${NC} branch...\n"
	cp Dockerfiles/Dockerfile.frontend frontend/Dockerfile
	cp Dockerfiles/nginx.conf.prod nginx/nginx.conf
	cp Dockerfiles/docker-compose.yml ./docker-compose.yml
	printf "${GREEN}Done!${NC}\n"
fi

######################################################################

#Rebuilding go (aka backend)
if [[ $1 == "restart" && $2 == "go" ]]
then
	printf "${GREEN}Rebuilding go container${NC}\n"
	docker-compose up --build -d go
fi

######################################################################

#Starting and stopping env

if [[ $1 == "up" && $2 == "" ]]
then
	printf "${GREEN}(Re)Building everything${NC}\n"
	docker-compose up --build -d
fi

if [[ $1 == "down" && $2 == "" ]]
then
	printf "${GREEN}Downing everything${NC}\n"
	docker-compose down
fi

######################################################################

#Help

if [[ $1 == "" && $2 == "" ]]
then
	printf "Usage: "
	printf "./devops.sh <arguments>\n"
	printf "\n\n"
	printf "Run without arguments to get this help\n"
	printf "${GREEN}up${NC} - build everything\n"
	printf "${GREEN}down${NC} - down everything\n"
	printf "${GREEN}switch dev${NC} - switch frontend to dev\n"
	printf "${GREEN}switch prod${NC} - switch frontend to prod\n"
	printf "${GREEN}restart go${NC} - rebuild go (AKA backend)\n"
fi
