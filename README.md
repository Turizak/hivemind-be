# Hivemind Backend Repo

[![Fly Deploy](https://github.com/rakazirut/hivemind-be/actions/workflows/fly.yml/badge.svg)](https://github.com/rakazirut/hivemind-be/actions/workflows/fly.yml)

## 🐝 What is Hivemind?

Hivemind is a reddit/forum style web application created by Tom Slanda and Robert Kazirut. This is the backend repo for that application. The codebase is written in Golang and leverages packages such as Gorm and Gin.

[Tom's Github Profile](https://github.com/slandath) | [Rob's Github Profile](https://github.com/rakazirut)

### 🐳 How to run locally via Docker

1. docker build -t hivemindbe .
2. docker run -p 8080:8080 --env-file ./.env hivemindbe
