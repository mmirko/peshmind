#!/usr/bin/env bash

set -euo pipefail

usage() {
	echo "Usage: $0 <save|restore> <profile>" >&2
	echo "       $0 clean [profile]" >&2
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROFILES_DIR="$SCRIPT_DIR/profiles"

if [[ $# -lt 1 ]]; then
	usage
	exit 1
fi

COMMAND="$1"
CURRENT_DIR="$PWD"

PROFILE_NAME=""
PROFILE_DIR=""

remove_privsw_files() {
	local target_dir="$1"
	find "$target_dir" -maxdepth 1 -type f -name 'privsw-*' -delete
}

copy_privsw_files() {
	local source_dir="$1"
	local target_dir="$2"

	while IFS= read -r -d '' file; do
		cp "$file" "$target_dir/"
	done < <(find "$source_dir" -maxdepth 1 -type f -name 'privsw-*' -print0)
}

case "$COMMAND" in
	save)
		if [[ $# -ne 2 ]]; then
			usage
			exit 1
		fi
		PROFILE_NAME="$2"
		PROFILE_DIR="$PROFILES_DIR/$PROFILE_NAME"
		if [[ ! -d "$PROFILE_DIR" ]]; then
			echo "Error: profile '$PROFILE_NAME' not found in '$PROFILES_DIR'." >&2
			exit 1
		fi
		remove_privsw_files "$PROFILE_DIR"
		copy_privsw_files "$CURRENT_DIR" "$PROFILE_DIR"
		;;
	restore)
		if [[ $# -ne 2 ]]; then
			usage
			exit 1
		fi
		PROFILE_NAME="$2"
		PROFILE_DIR="$PROFILES_DIR/$PROFILE_NAME"
		if [[ ! -d "$PROFILE_DIR" ]]; then
			echo "Error: profile '$PROFILE_NAME' not found in '$PROFILES_DIR'." >&2
			exit 1
		fi
		remove_privsw_files "$CURRENT_DIR"
		copy_privsw_files "$PROFILE_DIR" "$CURRENT_DIR"
		;;
	clean)
		if [[ $# -eq 1 ]]; then
			remove_privsw_files "$CURRENT_DIR"
		elif [[ $# -eq 2 ]]; then
			PROFILE_NAME="$2"
			PROFILE_DIR="$PROFILES_DIR/$PROFILE_NAME"
			if [[ ! -d "$PROFILE_DIR" ]]; then
				echo "Error: profile '$PROFILE_NAME' not found in '$PROFILES_DIR'." >&2
				exit 1
			fi
			remove_privsw_files "$PROFILE_DIR"
		else
			usage
			exit 1
		fi
		;;
	*)
		echo "Error: invalid command '$COMMAND'. Use 'save', 'restore', or 'clean'." >&2
		usage
		exit 1
		;;
esac

