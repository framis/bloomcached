Bloomcached
-----------

This is a bloom filter + Memcached. A bloom filter server with a simple TCP protocol similar to Memcached.

The bloom filter server is based on the work of https://github.com/willf/bloom. 

## Usage

```
go run server.go         				# Starts a default sever on port 3333
echo "ADD|hello" | nc localhost 3333    # Returns 201
echo "TEST|hello" | nc localhost 3333   # Returns 200|true
echo "TEST|new" | nc localhost 3333     # Returns 200|false
echo "ADD|new" | nc localhost 3333      # Returns 201
echo "TEST|new" | nc localhost 3333     # Returns 200|true
```

There is a simple client in Go, you can run the tests with `go test`


## Bloom Filter
See [here](https://www.jasondavies.com/bloomfilter/) for a nice explanation. Basically a bloom filter allows two operations, *add* and *test*. Test can tell for sure that an element is not in the bloomfilter. It can tell that the element may be in the filter. 

The main advantage of the bloomfilter is that it takes a very limited space. 

## Client
The client is a pretty simple TCP client, inspired from the [Memcached protocol](https://github.com/memcached/memcached/blob/master/doc/protocol.txt)

## Misc
This project has been built for fun and is not intended for production use.
