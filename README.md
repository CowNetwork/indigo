# indigo

Indigo is a microservice that enables you to bind roles and their respective permissions to user accounts.

It is running as an RPC service and uses the database to store and manage its data. Also a connection to a Kafka broker will be made, so that changes can be sent automatically to other services, if needed. For a view on what methods this microservice provides, have a look at [mooapis](https://github.com/CowNetwork/mooapis).

# Usage

First you have to be authenticated with the [Github Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry). You also need a running PostgreSQL instance and Kafka broker for it to work. See the [docker-compose.yml](https://github.com/CowNetwork/indigo/blob/main/docker-compose.yml) for a detailed example of how to get it all running together.

Then you can simply use the Docker image via:

```
docker run --rm ghcr.io/cownetwork/indigo:latest
```

Or you can use the compose file mentioned above and run:

```
docker-compose up -d
```

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
