# Golang Ex10 SDK

## Introduction

This is a simple golang interface around the extended exchange rust sdk [here](https://github.com/x10xchange/rust-crypto-lib-base). The SDK provides basic Account creation, Order submission and Market/FeeData fetching.

## Prerequisites

- Go 1.19 or later
- Rust toolchain (rustc, cargo)
- GCC or compatible C compiler
- Git

Note: Currently, the SDK is only compatible for use in linux-based x86_64 machines (including WSL)

## Project Structure

```
extended-sdk-golang/
├── README.md           # This file
└── src/
    ├── api_client.go      # REST API client for trading operations
    ├── base.go            # Base module with common HTTP functionality
    ├── config.go          # Configuration and domain models
    ├── markets.go         # Market data models
    ├── orders.go          # Order creation and management
    ├── sign.go            # Cryptographic signing with CGO bindings
    └── utils.go           # Utility functions
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

After building the Rust library, you must ensure to allow the go compiler to find the library by setting the library environment variable. Additionally, certain tests require testnet API keys and a private/public keypair.

```bash
export LD_LIBRARY_PATH="$(pwd):${LD_LIBRARY_PATH:-}"
export TEST_API_KEY=<your_api_key>
export TEST_PRIVATE_KEY=<your_private_key_hex>
export TEST_PUBLIC_KEY=<your_public_key_hex>
export TEST_VAULT=<your_vault_id>
```

Then you can run tests from the `src/` directory:

```bash
# Navigate to src directory
cd src

# Run all tests
go test -v

# Run specific test patterns
go test -run TestCreateOrderObject -v
go test -run TestGetOrderHash -v
go test -run TestSignMessage -v

# Run tests with race detection
go test -race -v
```

## Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/shopspring/decimal"
    sdk "github.com/extended-protocol/extended-sdk-golang/src"
)

func main() {
    // Create API client
    cfg := sdk.EndpointConfig{
        APIBaseURL: "https://api.starknet.sepolia.extended.exchange/api/v1",
    }
    
    // Initialize Stark account (example values - use your own)
    account, err := sdk.NewStarkPerpetualAccount(
        123,                                                           // vault
        "0x1234567890123456789012345678901234567890123456789012345678901234", // private key
        "0x0987654321098765432109876543210987654321098765432109876543210987", // public key
        "your-api-key-here",                                          // API key
    )
    if err != nil {
        log.Fatal("Failed to create account:", err)
    }
    
    client := sdk.NewAPIClient(cfg, account.APIKey(), account, 30*time.Second)
    defer client.Close()
    
    ctx := context.Background()
    
    // Get available markets
    markets, err := client.GetMarkets(ctx, []string{"BTC-USD"})
    if err != nil {
        log.Fatal("Failed to get markets:", err)
    }
    fmt.Printf("Available markets: %+v\n", markets)
    
    // Create an order
    if len(markets) > 0 {
        market := markets[0]
        
        // Setup Starknet domain (example for testnet)
        domain := sdk.StarknetDomain{
            Name:     "Perpetuals",
            Version:  "v0",
            ChainID:  "SN_SEPOLIA",
            Revision: "1",
        }
        
        // Order parameters
        nonce := 12345
        orderParams := sdk.CreateOrderObjectParams{
            Market:                   market,
            Account:                  *account,
            SyntheticAmount:          decimal.NewFromFloat(0.1),  // 0.1 BTC
            Price:                    decimal.NewFromFloat(50000), // $50,000
            Side:                     sdk.OrderSideBuy,
            Signer:                   account.Sign,
            StarknetDomain:          domain,
            PostOnly:                false,
            TimeInForce:             sdk.TimeInForceGTT,
            SelfTradeProtectionLevel: sdk.SelfTradeProtectionDisabled,
            Nonce:                   &nonce,
        }
        
        // Create order object
        order, err := sdk.CreateOrderObject(orderParams)
        if err != nil {
            log.Fatal("Failed to create order:", err)
        }
        
        // Submit order
        response, err := client.SubmitOrder(ctx, order)
        if err != nil {
            log.Fatal("Failed to submit order:", err)
        }
        
        fmt.Printf("Order submitted successfully: %+v\n", response)
    }
}
```

## Troubleshooting

### Build Issues
- **CGO errors**: Ensure you have a C compiler installed and the Rust library is built
- **Library not found**: Check that `liborderffi.so` exists in the root directory and `LD_LIBRARY_PATH` is set correctly
- **Runtime errors**: Verify the rpath is correctly set for your platform

### API Issues
- **400 Bad Request**: Check that all required fields are included in your order object and authentication headers are correct
- **Authentication errors**: Verify your API key is valid and properly set
- **Order validation errors**: Ensure order parameters (price, quantity, nonce) are within acceptable ranges

### Testing Issues
- **Missing environment variables**: Some tests require API keys and account credentials to be set
- **Library loading errors**: Ensure the Rust library is built and the `LD_LIBRARY_PATH` is correctly set before running tests

