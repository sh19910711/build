### Requirements

* ***Docker >= 1.8***

### Usage

#### Run Server

```
$ make run
```

#### Build Request

```
$ curl --form file=@example/app.tar --form callback=http://example.com/callback http://localhost:8080/builds
```

### Web API

#### GET /builds

Returns all builds:

```json
{
  "builds": [
    "Build Details"
  ]
}
```

#### POST /builds

Create new build:

```json
{
  "id": "Build ID"
}
```

#### GET /builds/:id

Returns details of the build:

```json
{
  "id": "Build ID",
  "jobs": [
    "Job Details"
  ]
}
```

#### GET /builds/:id/log.txt

Returns log messages of the build
