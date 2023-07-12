# Wisdom server

## Description

This is a Go project implementing a Wisdom server and client. The server provides a wisdom quote to clients after they solve a proof-of-work challenge. The client sends requests to the server to get these wisdom quotes.

## PoW explanation

There were no requirements mentioned for the PoW challenge (like ASIC-resistance or memory-hardness), so I decided to implement the simplest hashcash-like PoW algorithm using SHA-256 function. The difficulty of the challenge is adjusted by COMPLEXITY variable in .env file. The difficulty is the expected number of leading zeros in the hash of the concatenation of the challenge value and nonce, i.e. the challenge solution.

## Requirements

- Make
- Docker

## Usage

### Server

The server can be built and started with the following commands:

```bash
make build-server
make run-server
```

### Client

The client can be built and started with the following commands:

```bash
make build-client
make run-client
```
