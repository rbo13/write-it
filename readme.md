# write-it
[![Go Report Card](https://goreportcard.com/badge/github.com/rbo13/write-it)](https://goreportcard.com/report/github.com/rbo13/write-it)
[![Test Coverage](https://api.codeclimate.com/v1/badges/d6310e6cfc7b68ffb2bd/test_coverage)](https://codeclimate.com/github/rbo13/write-it/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/d6310e6cfc7b68ffb2bd/maintainability)](https://codeclimate.com/github/rbo13/write-it/maintainability)


> A sample dockerized application.


##### Building the application in docker
```sh
$ docker build -t <tag_name:latest> </path/to/entrypoint>
```


> Sample
```sh
$ docker build -t write-it:latest .
```


##### Running inside inside the docker container
```sh
$ docker run -it -p 1333:1333 write-it:latest
```