GRPC pool for groupcache
========================

A replacement for [groupcache](https://github.com/golang/groupcache) `HTTPPool` that uses GRPC to communicate with peers. 


Usage
-----

```go
server := grpc.NewServer()

p := NewGRPCPool("127.0.0.1:5000", server)
p.Set(peerAddrs...)

getter := groupcache.GetterFunc(func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
	dest.SetString(...)
	return nil
})

groupcache.NewGroup("grpcPool", 1<<20, getter)
lis, err := net.Listen("tcp", "127.0.0.1:5000")
if err != nil {
	log.Fatalf("Failed to start server")
}

server.Serve(lis)
```

Use `GRPCPoolOptions` to set the GRPC client dial options such as using compression, authentication etc. 
