#!/bin/bash

rm -rf docs build
docker run --rm --name slate2 -p 4567:4567 -v $(pwd):/srv/slate/ slate
mv build docs

