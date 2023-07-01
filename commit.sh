#!/bin/zsh

echo "commit message: "
read commit

git add .
git commit -m $commit