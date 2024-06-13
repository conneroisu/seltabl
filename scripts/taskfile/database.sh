#!/bin/bash
# file: makefile.database.sh
# url: https://github.com/dressingai/dressing-api/scripts/makefile.database.sh
# title: Combining SQL Schemas, Queries, and Seeds into Combined SQL Files (queries.sql, schema.sql, seeds.sql)

dbs=(
	"sqlite3"
)

# for each known database
for db in "${dbs[@]}"; do
	awk 'FNR==1{print ""}1' ./tools/data/"$db"/schemas/*.sql > "./tools/data/$db/combined/schema.sql"
	awk 'FNR==1{print ""}1' ./tools/data/"$db"/seeds/*.sql > "./tools/data/$db/combined/seeds.sql"
	awk 'FNR==1{print ""}1' ./tools/data/"$db"/queries/*.sql > "./tools/data/$db/combined/queries.sql"
done

for db in "${dbs[@]}"; do
	cd "./tools/data/$db" || echo "db $db not found" && exit
	echo "===== GENERATE $db ====="
	sqlc generate
	echo "^^^^^ GENERATE $db ^^^^^"
	cd "../../.." || echo "parent folder not found" && exit
	rm "./tools/data/$db/db.go"
done
