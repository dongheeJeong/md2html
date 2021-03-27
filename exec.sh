#!/bin/sh

go run parser.go index.md > index.html

cat misc/pre.html index.html misc/post.html > article.html

cat article.html
brave article.html

