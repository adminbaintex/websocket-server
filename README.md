# Websocket server

Run
====

Run example server:

    go run example_server.go -addr localhost:8081 -path api

Run example clients:

    go run example_client.go -addr localhost:8081 -path api
    go run example_client.go -addr localhost:8081 -path api

To view server heap:

    http://localhost:6060/debug/pprof/heap

To view server running goroutines:

    http://localhost:6060/debug/pprof/goroutine?debug=2

Dependencies 
=============

Dependencies are managed with govendor.

* Initialize "vendor" directory
```
govendor init
```

* List packages used in the application
```
govendor list
```

* Add external packages in GOPATH to vendor folder
```
govendor add +external
```
