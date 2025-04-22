#!/bin/bash

cat privsw-* | grep "switch(" > data.pl
cat privsw-* | grep "switchname(" >> data.pl
cat privsw-* | grep "seen(" >> data.pl
cat rules.pl >> data.pl
cat server.pl >> data.pl
