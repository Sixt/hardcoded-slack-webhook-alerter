#!/bin/bash

# Put a list of repos to clone here
services=(
	"MyRepoOne"
	"SomeOtherRepo"
	"com.service.example"
)

# Replace this with your GitHub organization/user name
org="MyOrg"

fetchService() {
	printf "\n###### Fetching $1 #####\n"

	#clone the repo if it's not already there
	[ ! -d $1 ] && git clone git@github.com:$org/$1.git 

	echo ""
}

for service in "${services[@]}"; do
	fetchService $service
done
