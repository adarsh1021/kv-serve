# KVServe

![Go](https://img.shields.io/github/go-mod/go-version/adarsh1021/kv-serve?logo=go)
![GitHub](https://img.shields.io/github/license/adarsh1021/kv-serve)

KVServe is a persistent key-value database as a service based on [LevelDB](https://github.com/google/leveldb).

## Features

- Uses LevelDB as the backend for the key-value database.
- Fast and simple HTTP API.
- Read/write across multiple databases without any performance impact.

## Usage

To build:

```
go build -o kv-serve
```

Options:

- --data-dir (-d)  
  The path to store LevelDB files.
- --max-db-cache-entries (-c)  
  The maximum number of open db pointers at a time.
- --port (-p)  
  The port to run the HTTP server on.

Create a new key-value database `my-db`:

```
curl -X POST http://localhost:9090/db/my-db
my-db created
```

Set a new key `hello` in `my-db`:

```
curl -d "world" -X POST http://localhost:9090/db/my-db/hello
ok
```

Get the value of `hello` from `my-db`:

```
curl -X GET http://localhost:9090/db/my-db/hello
world
```
