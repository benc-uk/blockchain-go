# Blockchain Learning in Go

A trivial blockchain with a REST API

- Based on Proof of Work PoW hashing
- No persistence
- Doesn't really do anything

# Blockchain Implementation

See [blockchain/base.go](./blockchain/base.go)

# Running Locally

Use the makefile :)
```
make run
```

The API server will start on port 8080 (or whatever `PORT` env var is set to)

```text
help                 💬 This help message :)
install-tools        🔮 Install dev tools into project bin directory
lint                 🌟 Lint & format, will not fix but sets exit code on error
lint-fix             🔍 Lint & format, will try to fix errors and modify code
run                  🏃 Run application, used for local development
```

# API

The API is RESTful and supports the following operations

See [blockchain.http](./blockchain.http) for example of calling the API and sample requests

| Method | Path                     | Description             | Body          | Returns            |
| ------ | ------------------------ | ----------------------- | ------------- | ------------------ |
| GET    | /chain/list              | Dump the whole chain    | None          | Array of _Block_   |
| GET    | /chain/validate          | Check chain integrity   | None          | Status of chain    |
| GET    | /block/_{hash}_          | Get a single block      | None          | _Block_            |
| POST   | /block                   | Add a transaction block | _Transaction_ | Hash of new block  |
| PUT    | /block/tamper/_{hash}_   | Tamper with block data  | None          | Status             |
| GET    | /block/validate/_{hash}_ | Check block integrity   | None          | Integrity of block |

```go
type Block struct {
	Timestamp    time.Time
	Hash         string
	PreviousHash string
	Data         string
	Nonce        int
	Index        int
}

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}
```
