#!/bin/bash

source "$(pwd)"/tests/env.sh
docker-compose -f "$(pwd)"/tests/docker-compose.yml down
unset TF_ACC MONGODB_URL MONGODB_USERNAME MONGODB_PASSWORD
