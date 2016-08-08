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
