#!/bin/bash


if [ "$#" -ne 1 ]; then
    echo "Useage: ./import.sh softradius"
		exit 1
fi

DBNAME=$( echo $1 | sed s/\\///g )

cd $DBNAME

echo "Droping... $DBNAME"

recli "r.dbDrop(\"$DBNAME\")"


for f in *
do
  table=$(echo $f | sed "s/\(.*\)_[0-9]\+.json/\1/g")
  table=$(echo $table | sed "s/-/_/g")
	table=$(echo $table | sed s/".json"//g )
	echo $table
 
  echo "rethinkdb import -f $f --table $DBNAME.$table --force"
  rethinkdb import -f $f --table $DBNAME.$table --force

done

cd ..

