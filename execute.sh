#!/usr/bin/env bash

PROFILES_DIR="$(pwd)/kbpool/profiles"

set -euo pipefail

usage() {
	echo "Usage: $0 clean" >&2
	echo "       $0 sim <arg>" >&2
	echo "       $0 build <arg>" >&2
	echo "       $0 serve" >&2
	echo "       $0 save <arg>" >&2
	echo "       $0 restore <arg>" >&2
	echo "       $0 show <arg>" >&2
	echo
}

clean() {
	rm -f peshmind.json
	rm -f peshmind
	rm -f simout.dot
	rm -f kbpool/data.pl
	cd kbpool >/dev/null 2>&1
	./profile.sh clean
	cd - >/dev/null 2>&1
	:
}

sim() {
	local arg="$1"
	clean
	go build
	cp simpool/peshmind.json .
	./peshmind simulate -config peshmind.json -d simout.dot -o kbpool/privsw-sim.pl -s "$arg" 
	:
}

build() {
	local arg="$1"
	clean
	restore "$arg"
	go build
	:
}

serve() {
	go build
	./peshmind serve -config peshmind.json
	:
}

save() {
	local arg="$1"
	if [[ ! -d "$PROFILES_DIR" ]]; then
		echo "Error: profiles directory '$PROFILES_DIR' not found." >&2
		exit 1
	fi
	if [[ -d "$PROFILES_DIR/$arg" ]]; then
		echo "profile '$arg' already exists. Do you want to overwrite it? (y/N)" >&2
		read -r answer
		if [[ "$answer" != "y" ]]; then
			echo "Aborting." >&2
			exit 1
		fi
		rm -rf "$PROFILES_DIR/$arg"
	fi

	mkdir -p "$PROFILES_DIR/$arg"
	cp -a peshmind.json "$PROFILES_DIR/$arg/"
	cd kbpool >/dev/null 2>&1
	./profile.sh save "$arg"
	cd - >/dev/null 2>&1
	:
}

show() {
	local arg="$1"
	case "$arg" in
		"simdot")
			if [[ ! -f simout.dot ]]; then
				echo "Error: simout.dot not found. Please run '$0 sim <arg>' first." >&2
				exit 1
			fi
			cat simout.dot | fdp -Txlib
			;;
		"inferreddot")
			go build
			./peshmind dot -config peshmind.json | fdp -Txlib
			;;
		*)
			echo "Error: invalid argument '$arg' for show command." >&2
			exit 1
			;;
	esac
}		

restore() {
	local arg="$1"
	if [[ ! -d "$PROFILES_DIR/$arg" ]]; then
		echo "Error: profile '$arg' not found in '$PROFILES_DIR'." >&2
		exit 1
	fi
	cp -a "$PROFILES_DIR/$arg/peshmind.json" .
	cd kbpool >/dev/null 2>&1
	./profile.sh restore "$arg"
	cd - >/dev/null 2>&1
	:
}

if [[ $# -lt 1 ]]; then
	usage
	exit 1
fi

command="$1"

case "$command" in
	clean)
		if [[ $# -ne 1 ]]; then
			usage
			exit 1
		fi
		clean
		;;
	serve)
		if [[ $# -ne 1 ]]; then
			usage
			exit 1
		fi
		serve
		;;
	sim)
		if [[ $# -ne 2 ]]; then
			usage
			exit 1
		fi
		sim "$2"
		;;
	build)
		if [[ $# -ne 2 ]]; then
			usage
			exit 1
		fi
		build "$2"
		;;
	save)
		if [[ $# -ne 2 ]]; then
			usage
			exit 1
		fi
		save "$2"
		;;
	show)
		if [[ $# -ne 2 ]]; then
			usage
			exit 1
		fi
		show "$2"
		;;
	restore)
		if [[ $# -ne 2 ]]; then
			usage
			exit 1
		fi
		restore "$2"
		;;
	*)
		echo "Error: invalid command '$command'." >&2
		usage
		exit 1
		;;
esac
