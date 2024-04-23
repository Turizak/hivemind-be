# hivemind-be

[![Fly Deploy](https://github.com/rakazirut/hivemind-be/actions/workflows/fly.yml/badge.svg)](https://github.com/rakazirut/hivemind-be/actions/workflows/fly.yml)

Hivemind Backend Repo

Docker:
1. docker build -t hivemindbe .
2. docker run -p 8080:8080 --env-file ./.env hivemindbe