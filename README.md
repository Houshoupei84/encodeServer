# encodeServer

## Synopsis
The encodeServer project provides a password hashing service over http.
It maintains the processed hashes for the duration of the process execution time.

## Code Example
This example starts a new encodeServer listening on port 123:
``` {.sourceCode .golang}
s := server.Server{"123"}
s.Run()
```

## Usage
To use as a standalone application:
go run main.go [port]
(or)
go build main.go && ./main [port]

If port is not specified, port 8080 will be used as a default.

## API Reference
When an encodeServer is running, it will process the following http requests:

/hash POST password=example
Responds with the id for the hash.
After a 5 second delay, computes the base64-encoded RSA512 hash of the given password and stores it.

/hash/N GET
Responds with the hash corresponding to N, where N is a hash id.

/stats GET
Return a summary of the total number of requests and average response time in microseconds.

/shutdown GET
Gracefully shutdown the server once existing requests have completed.
