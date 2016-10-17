## Build server

The build server based on docker :whale:

### Why

* To hide docker hosts from web
* To provide cross build environment with docker

### Requirements

* ***Docker >= 1.8***

### Usage

#### Run server

```
$ make build
$ ./app
```

#### Create a build request

Send a json text to `POST /builds`

```
{"id":"here-is-build-id"}
```
