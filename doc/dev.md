```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/easybackup.proto
```

```bash
go run cmd/* server run --addr 0.0.0.0:50051 --id 1  --raft_addr 127.0.0.1:5000
go run cmd/* server run --addr 0.0.0.0:50061 --id 2  --raft_addr 127.0.0.1:6000
go run cmd/* server run --addr 0.0.0.0:50071 --id 3  --raft_addr 127.0.0.1:7000
go run cmd/* server run --addr 0.0.0.0:50081 --id 4  --raft_addr 127.0.0.1:8000


go run cmd/* server addnode --addr 0.0.0.0:50051 --id 2  --raft_addr 127.0.0.1:6000
go run cmd/* server addnode --addr 0.0.0.0:50051 --id 3  --raft_addr 127.0.0.1:7000
go run cmd/* server addnode --addr 0.0.0.0:50051 --id 4  --raft_addr 127.0.0.1:8000
```
