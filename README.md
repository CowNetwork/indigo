# indigo

Indigo is a microservice that enables you to bind roles and their respective permissions to user accounts.

# Environment Variables

| Variable | Default | Description |
| -------- | ------- | ----------- |
| `INDIGO_SERVICE_POSTGRES_URL` | `localhost:5432` | The url the postgres listens to. |
| `INDIGO_SERVICE_POSTGRES_USER` | `test` | The user to connect to the postgres. |
| `INDIGO_SERVICE_POSTGRES_PASSWORD` | `password` | The password to connect to the postgres. |
| `INDIGO_SERVICE_POSTGRES_DB` | `test` | The database to connect to the postgres. |
| `INDIGO_SERVICE_POSTGRES_SCHEMA` | `public` | The schema to connect to. |
| `INDIGO_SERVICE_HOST` | `localhost` | Host to bind the service to. |
| `INDIGO_SERVICE_PORT` | `6969` | Port to bind the service to. |
| `INDIGO_SERVICE_KAFKA_BROKERS` | `127.0.0.1:9092` | Kafka brokers to connect to. |
| `INDIGO_SERVICE_KAFKA_TOPIC` | `cow.global.indigo` | Kafka topic to send events to. |
| `INDIGO_SERVICE_CLOUDEVENTS_SOURCE` | `cow.global.indigo-service` | CloudEvents source uri. |