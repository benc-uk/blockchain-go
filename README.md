# 🏋️ Go Blockchain 

A simple blockchain written in Go, with a REST API and command line client

- Based on *Proof of Work* SHA hashing
- Data persistence of the blockchain using [BoltDB](https://github.com/etcd-io/bbolt) and [gob encoding](https://pkg.go.dev/encoding/gob)
- REST API using [mux](https://github.com/gorilla/mux)
- Command line client using [Cobra](https://github.com/spf13/cobra)
- Mock/simple transaction system

# 🔗 Blockchain Implementation

See:

- [blockchain/block.go](./blockchain/block.go)
- [blockchain/chain.go](./blockchain/chain.go)

# 🏃‍♂️ Usage & Running Locally

First build the client

```bash
make build-client
```

Then run the API server/backend

```bash
make run-server
```

In a separate terminal run the client CLI binary

```bash
./out/blockchain
```

The command line supports several commands, such as `add`, `list`, `get`

```text
Usage:
  blockchain [command]

Available Commands:
  add         Add a transaction to the blockchain
  completion  Generate the autocompletion script for the specified shell
  get         Get a single block by its hash
  help        Help about any command
  list        List the whole block chain
  validate    Check integrity
```

To create a transaction, mine the block and add it to the blockchain, run `blockchain add` with a sender name, recipient and amount value.

```bash
./out/blockchain add -s "David Bowie" -r "Lou Reed" -a 39.66
```

Verify with `blockchain list` or `blockchain get`

Further makefile commands are available

```text
help                 💬 This help message :)
get-tools            🔮 Install dev tools into project bin directory
lint                 🌟 Lint & format, will not fix but sets exit code on error
lint-fix             🔍 Lint & format, will try to fix errors and modify code
run-server           🏃 Run server API locally, with hot reload
build-server         🔨 Build server API, resulting binary in out/blockchain-api
build-client         🔨 Build client command line tool, resulting binary in out/blockchain
clean                💣 Clean up, database and temp files
```

# 🤖 Configuration

### Server

- `PORT` - Set the port the server will listen on, 8080 is the default
- `CHAIN_DIFFICULTY` - Set the difficulty level, defaults to '5', note that values higher than '6' are not advised, unless you have a super computer!

### Client

- `API_ENDPOINT` - Configure where the server is running, defaults to `https://localhost:8080`


# 🌐 API Reference

The API is RESTful and supports the following operations

See [blockchain.http](./blockchain.http) for example of calling the API and sample requests

| Method | Path                     | Description            | Body    | Returns            |
| ------ | ------------------------ | ---------------------- | ------- | ------------------ |
| GET    | /chain/list              | Dump the whole chain   | None    | Array of _Block_   |
| GET    | /chain/validate          | Check chain integrity  | None    | Status of chain    |
| GET    | /block/_{hash}_          | Get a single block     | None    | _Block_            |
| POST   | /block                   | Add a mined block      | _Block_ | _Block_            |
| PUT    | /block/tamper/_{hash}_   | Tamper with block data | None    | Status             |
| GET    | /block/validate/_{hash}_ | Check block integrity  | None    | Integrity of block |

```go
type Block struct {
	Timestamp    time.Time
	Hash         string
	PreviousHash string
	Data         string
	Nonce        int
}

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}
```