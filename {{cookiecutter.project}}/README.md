# {{cookiecutter.project_name}} [Backend Application] ![GO][go-badge]

[go-badge]: https://img.shields.io/badge/Go-v1.19-blue

## Build & Run (Locally) 🙌👨‍💻🚀
### Prerequisites
- go 1.19
- docker & docker-compose
- [golangci-lint](https://github.com/golangci/golangci-lint) (<i>optional</i>, used to run code checks)
- [swag](https://github.com/swaggo/swag) (<i>optional</i>, used to re-generate swagger documentation)

Use `make swag` to generate docs(required), make local` to run docker compose, `make run` to run the project, `make lint` to check code with linter.

### Hot reload
 
First step install nodemon like this -> npm i nodemon -g

Run `nodemon -x go run cmd/main.go --signal SIGTERM -e go --verbose`

### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3005
