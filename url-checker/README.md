
# Concurrent URL Checker

A lightweight concurrent URL checker built with Go.

## Features

- Concurrent HTTP requests
- Goroutine-based worker pool
- Job distribution using channels
- Result collection using buffered channels
- WaitGroup synchronization
- HTTP request timeouts
- Request duration measurement

## Tech Stack

- Go
- net/http
- sync
- time

## Run Locally

```bash
go run .
```

## Run with Race Detector

```bash
go run -race .
```

## Concepts Practiced

- Goroutines
- Channels
- Buffered channels
- Worker pools
- WaitGroups
- HTTP clients
- Error handling

## Status

Learning project. Built to practise Go concurrency.
