#!/bin/bash

source "$(pwd)"/tests/env.sh
docker-compose -f "$(pwd)"/tests/docker-compose.yml up -d
"$(pwd)"/tests/wait-mongodb-docker.sh "$(pwd)"/tests/docker-compose.yml
