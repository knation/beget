Go web service for producing to a Kafka topic over HTTP.

**NOT FOR PRODUCTION USE**
- Needs testing with live Kafka cluster, both with a single node and multiple brokers
- Support for more Kafka connection options
- Implement more robust configuration input via `viper`

# Motivation
In order to produce to a Kafka via HTTP, you need a proxy. You could use the [Confluent Rest Proxy](https://github.com/confluentinc/kafka-rest), but it can be difficult to configure and deploy. Also, if you want to transform/validate data or do anything else before producing, you'd need to have another service in front of it anyway.

If hosted on Confluent, you can use their provided REST proxy, but then you're subject to their limits and lack of observability/flexibility. If you want to do any sort of validation, you'd still have to stand up another service in between or move your business logic to a layer upstream.

The idea with this project is to have a simple web service that does nothing but accept HTTP requests and produces the payload to Kafka. Though, as a simple Go service, you can extend it if you'd like to include validation, transformation, or anything else.

## Use cases

* Use behind a reverse proxy to produce valid messages. For example, validate requests in nginx and create a subrequest to `beget` if the message is okay to put into Kafka.
* Use in tandem with serverless code that is unable to implement a Kafka driver. For example, producing to Kafka via a Cloudflare Worker.
* When you have multiple microservices that all need to produce to Kafka, it may be easier to have them make an HTTP request to a common shared microservice rather than needed to fully implement a driver, especially if they're written in different languages.
* You may have a need to scale your Kafka message production independently from other microservices.

# Usage

The project entrypoint may be found at `cmd/main.go`. You can run locally with: `TOPICS=events go run cmd/main.go`.

## Configuration

beget uses [`viper`](https://github.com/spf13/viper) for managing configuration. This means that the configuration options below can also be provided via flag or ENV. Possible configuration options are:

```yaml
app:
  mode: debug # Application run mode (debug|release). Default: debug

server:
  port: 8080 # Web service port. Default: 8080

kafka:
  brokers: # List of kafka brokers to connect o 
    - broker1
    - broker2
    - ...
  topics: # List of kafka topics to allow
    - topic1
    - topic2
    - ...
```

In "debug" mode, the service does not connect to Kafka and messages are just logged.


## Producing to a topic
To produce to a topic, make a `POST` request to `/produce`, passing the topic and message as JSON in the body. For example, to produce a simple message (`{"foo":"bar"}`) to the "events" topic:
```
curl --request POST 'http://localhost:8080/produce' \
     --header 'Content-Type: application/json' \
     --data-raw '{"topic":"events","message":{"foo":"bar"}}'
```

### Body Parameters
| Parameter | Description                          | Required | Default |
|-----------|--------------------------------------|----------|---------|
| `topic`   | The topic to produce the message to. | Yes      |         |
| `value` | The message to produce to the topic. | Yes      |         |
| `key` | The message key. | No      |         |

## Health check
The service will respond with a 200 status code on any request to `/healthz`.

# Deploying

Initially, this project is intended for deployment on services like [Heroku](https://www.heroku.com/) or [Render](https://render.com/). Once stable, a Dockerfile will be created for containerized deployment.

# Contributing

* Please file an issue if there's a bug or feature request.
* Pull requests are welcome and will be reviewed/merged as appropriate.

# License

MIT License
