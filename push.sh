#! /bin/bash -

dataTime=$(date "+%F %T")
# add 
git commit -a -m "$dataTime"

git push
