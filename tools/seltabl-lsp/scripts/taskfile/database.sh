#!/bin/bash

dbs=(
	"master"
)

# for each known database
for db in "${dbs[@]}"; do
	awk 'FNR==1{print ""}1' ./data/"$db"/schemas/*.sql > "./data/$db/combined/schema.sql"
	awk 'FNR==1{print ""}1' ./data/"$db"/seeds/*.sql > "./data/$db/combined/seeds.sql"
	awk 'FNR==1{print ""}1' ./data/"$db"/queries/*.sql > "./data/$db/combined/queries.sql"
done

for db in "${dbs[@]}"; do
	cd "./data/$db" || echo "db $db not found" && exit
	echo "===== GENERATE $db ====="
	sqlc generate
	echo "^^^^^ GENERATE $db ^^^^^"
	cd "../.." || echo "parent folder not found" && exit
	rm "./data/$db/db.go"
done
