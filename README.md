# Golang Ex10 SDK

## Introduction

This is a simple golang interface around the extended exchange rust sdk [here](https://github.com/x10xchange/rust-crypto-lib-base). The SDK provides Go bindings for cryptographic operations including order hash generation and message signing for the extended exchange protocol.

## Features

- **Order Hash Generation**: Compute order hashes for trading operations
- **Message Signing**: Sign messages using private keys with ECDSA
- **Cross-platform**: Built with CGO bindings to Rust library

## Prerequisites

- Go 1.19 or later
- Rust toolchain (rustc, cargo)
- GCC or compatible C compiler
- Git

## Project Structure

```
golang-ex10-sdk/
├── README.md           # This file
├── sdk.go             # Main SDK implementation with CGO bindings
├── sdk_test.go        # Unit tests
└── rust-lib/          # Rust library source code
    └── target/
        └── release/   # Built Rust library (.so file)
```

## Building the Rust Library

1. Navigate to the rust-lib directory:
```bash
cd rust-lib
```

2. Build the release version of the Rust library:
```bash
cargo build --release
```

This will generate the shared library at `rust-lib/target/release/liborderffi.so` (Linux) or equivalent for your platform. You must then copy the library to the source directory.

Alternatively, run `build-lib.sh` in the root directory.

## Running Tests

After building the Rust library, you must ensure to allow the go compiler to find the library by setting the library environment variable
```bash
export LD_LIBRARY_PATH="$(pwd):${LD_LIBRARY_PATH:-}"
```

```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run specific test
go test -run TestGoGetOrderHash
```

## Usage Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/extended-protocol/extended-sdk-golang"
)

func main() {
    // Generate order hash
    hash, err := sdk.GetOrderHash(
        "100", "0x2", "100",
        "0x1", "-156",
        "0x1", "74",
        "100", "123",
        "0x5d05989e9302dcebc74e241001e3e3ac3f4402ccf2f8e6f74b034b07ad6a904",
        "Perpetuals", "v0", "SN_SEPOLIA", "1",
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Order hash:", hash)
    
    // Sign a message
    sig, err := sdk.SignMessage(hash, "0x1234def56789012345678901234567890123456789012345678901234567890")
	if err != nil {
		log.Fatal("SignMessage failed: %v", err)
	}
    fmt.Println("Signature:", signature)
}
```

## Troubleshooting

- **CGO errors**: Ensure you have a C compiler installed and the Rust library is built
- **Library not found**: Check that `liborderffi.so` exists in the root directory
- **Runtime errors**: Verify the rpath is correctly set for your platform

