#!/bin/bash

# Put a list of paths here where you keep your git repos
directories=(
    "/Users/me/go/"
    "/Users/me/java"
    "/Users/me/php/src"
)

gitPullAll() {
    for d in ./*/ ; do 
        echo "Updating $d..."
        cd "$d"

        #pull latest master/main
        git clean -df
        BRANCH=`git branch -l master main | sed -r 's/^[* ] //' | head -n 1`
        git checkout $BRANCH
        git checkout .
        git pull

        cd ..
    done
}

for dir in "${directories[@]}"; do
	cd $dir
    gitPullAll
done
