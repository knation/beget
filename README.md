**STATUS: This project is currently in development.**

---

Go web service for producing to a Kafka topic over HTTP.

# Motivation
In order to produce to a Kafka via HTTP, you need a proxy. You could use the [Confluent Rest Proxy](https://github.com/confluentinc/kafka-rest), but it can be difficult to configure and deploy. Also, if you want to transform/validate data or do anything else before producing, you'd need to have another service in front of it anyway.

The idea with this project is to have a simple web service that does nothing but accept HTTP requests and publishes the payload to Kafka. Though, as a simple Go service, you can extend it if you'd like to include validation, transformation, or anything else.

## Use cases

* Use behind a reverse proxy to produce valid messages. For example, validate requests in nginx and create a subrequest to `beget` if the message is okay to put into Kafka.
* Use in tandem with serverless code that is unable to implement a Kafka driver. For example, producing to Kafka via a Cloudflare Worker.
* When you have multiple microservices that all need to produce to Kafka, it may be easier to have them make an HTTP request to a common shared microservice rather than needed to fully implement a driver, especially if they're written in different languages.
* You may have a need to scale your Kafka message production independently from other microservices.

# Deploying

Initially, this project is intended for deployment on services like [Heroku](https://www.heroku.com/) or [Render](https://render.com/). Once stable, a Dockerfile will be created for containerized deployment.

# Contributing

* Please file an issue if there's a bug or feature request.
* Pull requests are welcome and will be reviewed/merged as appropriate.

# License

MIT License
