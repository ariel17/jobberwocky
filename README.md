# Jobberwocky home challenge

Hi there! This is my implementation for the Jobberwocky's home challenge. [Read
the requirements here](./docs/requirements.pdf).

## How to execute it

Compose is used to build the infrastructure required for
local execution and the Avature's external source.

```bash
# console 1
$ docker-compose up  # use -d to detach console

# console 2
$ curl -v "http://localhost:8080/search?title=java"
```

Once up and running, read the [Swagger 
documentation](http://localhost:8080/docs) for details
on API usage.

## Architecture diagram

![architecture diagram](./docs/architecture.png)

## Sequence diagram

![sequence diagram](./docs/sequence.png)

## Models diagram

![models diagram](./docs/models.png)