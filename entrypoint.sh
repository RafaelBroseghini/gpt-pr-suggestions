#!/bin/bash -l

export GITHUB_TOKEN=$1
export OPEN_AI_TOKEN=$2
export OPEN_AI_MODEL=$3
export IGNORED_FILE_EXTENSIONS=$4

cd /
mv ./gpt-pr-suggestions $GITHUB_WORKSPACE
cd $GITHUB_WORKSPACE
./gpt-pr-suggestions
