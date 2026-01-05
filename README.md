# Assignment 2 – Concurrent Task Service

This project implements a simple(for someone but not for me) concurrent task processing service in Go.

Initially,I selected **Option 2 (Generic Task Queue Wrapper)**, but during development it turned out to be unnecessarily complex due to generics and channel interactions.  
Therefore, the implementation was refactored to **Option 1 (Generic In-Memory Repository)**, while task processing is still handled concurrently using Go channels and worker goroutines.

## Features
- Concurrent task processing with goroutines
- Buffered channel used as a task queue
- Generic in-memory repository (Option 1)
- Thread-safe storage using mutexes
- Simple(FOR ME WAS NOT SIMPLE) HTTP API

## Run
go run .

## Testing
So for testing I used curl as it was required
And when it worked I was jumping, flying you can not imagine how happy I was.
curl.exe -X POST -d "payload=RassulDalban" http://localhost:8080/tasks
curl.exe http://localhost:8080/tasks
curl.exe http://localhost:8080/tasks/1
curl.exe http://localhost:8080/stats

## Conclusion

Thank you Nursultan Agai!!!