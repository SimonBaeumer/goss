# Random development scripts

## debian

`Vagrantfile` to develop on `debian` based systems.

```bash
vagrant up
```

## ssl

Create a https server with client certificate authentication

```bash
./generate-certificates.sh
go run https_server.go
``` 

## http

Start a http server for testing http request, i.e. headers.

```bash
go run http_test_server.go
```