#!/bin/bash

switchs=$(../peshmind list -r --config ../peshmind.json)

rm -f data.pl
touch data.pl

for file in $switchs; do
    if [ -f "$file.pl" ]; then
	cat "$file.pl" | grep "switch(" >> data.pl
    fi
done

for file in $switchs; do
    if [ -f "$file.pl" ]; then
	cat "$file.pl" | grep "switchname(" >> data.pl
    fi
done

for file in $switchs; do
    if [ -f "$file.pl" ]; then
	cat "$file.pl" | grep "seen(" >> data.pl
    fi
done

cat rules.pl >> data.pl
cat server.pl >> data.pl
