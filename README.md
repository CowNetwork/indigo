# indigo

Indigo is a microservice that enables you to bind roles and their respective permissions to user accounts.

# Environment Variables

| Variable | Default | Description |
| -------- | ------- | ----------- |
| `POSTGRES_URL` | `localhost:5432` | The url the postgres listens to. |
| `POSTGRES_USER` | `test` | The user to connect to the postgres. |
| `POSTGRES_PASSWORD` | `password` | The password to connect to the postgres. |
| `POSTGRES_DB` | `test` | The database to connect to the postgres. |
| `INDIGO_SERVICE_HOST` | `localhost` | Host to bind the service to. |
| `INDIGO_SERVICE_PORT` | `6969` | Port to bind the service to. |