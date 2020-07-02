# go-eth2-wallet-store-vault

[![Tag](https://img.shields.io/github/tag/Stakedllc/go-eth2-wallet-store-vault.svg)](https://github.com/Stakedllc/go-eth2-wallet-store-vault/releases/)
[![License](https://img.shields.io/github/license/Stakedllc/go-eth2-wallet-store-vault.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/Stakedllc/go-eth2-wallet-store-vault?status.svg)](https://godoc.org/github.com/Stakedllc/go-eth2-wallet-store-vault)
[![Travis CI](https://img.shields.io/travis/Stakedllc/go-eth2-wallet-store-vault.svg)](https://travis-ci.org/Stakedllc/go-eth2-wallet-store-vault)
[![codecov.io](https://img.shields.io/codecov/c/github/Stakedllc/go-eth2-wallet-store-vault.svg)](https://codecov.io/github/Stakedllc/go-eth2-wallet-store-vault)
[![Go Report Card](https://goreportcard.com/badge/github.com/Stakedllc/go-eth2-wallet-store-vault)](https://goreportcard.com/report/github.com/Stakedllc/go-eth2-wallet-store-vault)

Hashicorp Vault-based store for the [Ethereum 2 wallet](https://github.com/wealdtech/go-eth2-wallet).


## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-eth2-wallet-store-vault` is a standard Go module which can be installed with:

```sh
go get github.com/wealdtech/go-eth2-wallet-store-vault
```

## Usage

In normal operation this module should not be used directly.  Instead, it should be configured to be used as part of [go-eth2-wallet](https://github.com/wealdtech/go-eth2-wallet).

The Vault store has the following options:

  - `id`: an ID that is used to differentiate multiple stores created by the same account.  If this is not configured an empty ID is used
  - `passphrase`: a key used to encrypt all data written to the store.  If this is not configured data is written to the store unencrypted (although wallet- and account-specific private information may be protected by their own passphrases)

### Example

```go
package main

import (
	e2wallet "github.com/wealdtech/go-eth2-wallet"
	vault "github.com/Stakedllc/go-eth2-wallet-store-vault"
)

func main() {
    // Set up and use an encrypted store
    store, err := vault.New(vault.WithPassphrase([]byte("my secret")))
    if err != nil {
        panic(err)
    }
    e2wallet.UseStore(store)

    // Set up and use an encrypted store with a non-default vault address
    store, err = vault.New(vault.WithPassphrase([]byte("my secret")), vault.WithVaultAddress("https://my-secret-vault-server"))
    if err != nil {
        panic(err)
    }
    e2wallet.UseStore(store)

    // Set up and use an encrypted store with a different vault role
    store, err = vault.New(vault.WithPassphrase([]byte("my secret")), vault.WithRole("eth2role"))
    if err != nil {
        panic(err)
    }
    e2wallet.UseStore(store)

    // Set up and use an encrypted store with data stored in a different part of vault
    store, err = vault.New(vault.WithPassphrase([]byte("my secret")), vault.WithVaultSubPath("eth-secrets"))
    if err != nil {
        panic(err)
    }
    e2wallet.UseStore(store)
}
```

## Maintainers

Max Bucci: [@mbucci](https://github.com/mbucci).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/Stakedllc/go-eth2-wallet-store-vault/issues).
