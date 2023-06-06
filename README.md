# IBAX Client
[![Go Reference](https://pkg.go.dev/badge/github.com/IBAX-io/ibax-cli.svg)](https://pkg.go.dev/github.com/IBAX-io/ibax-cli)

IBAX official command line tool.

It is implemented based on [IBAX-SDK](https://github.com/IBAX-io/go-ibax-sdk), accessing IBAX network through rpc service


### Build from Source
Building `ibax-cli` requires both a Go (version 1.18 or later) and a C compiler.

```shell
make all
```

## Command

### config
The config command is used to generate a default configuration file

### console
The console command starts the console program, integrates most commands, and includes auto-completion functions


## Run

1. Create the node configuration file:
```bash
$  ibax-cli config
```
The `config` command contains flags: `connect,cryptoer,dataDir,ecosystem,hasher string,keysDir,port`.Can be viewed through help naming

2. start console or run command
```bash
    ibax-cli console
    ibax-cli version
```
