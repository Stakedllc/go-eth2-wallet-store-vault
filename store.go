// Copyright 2019, 2020 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vault

import (
	"io/ioutil"

	"github.com/hashicorp/vault/api"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// options are the options for the S3 store
type options struct {
	passphrase   []byte
	role         string
	vaultAddress string
	vaultSubPath string
}

// Option gives options to New
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithVaultAddress sets the vault address to connect to for the store
func WithVaultAddress(vaultAddress string) Option {
	return optionFunc(func(o *options) {
		o.vaultAddress = vaultAddress
	})
}

// WithPassphrase sets the passphrase for the store.
func WithPassphrase(passphrase []byte) Option {
	return optionFunc(func(o *options) {
		o.passphrase = passphrase
	})
}

// WithRole sets the role for the store.
func WithRole(role string) Option {
	return optionFunc(func(o *options) {
		o.role = role
	})
}

// WithVaultSubPath sets thewallet name for the Store
func WithVaultSubPath(vaultSubPath string) Option {
	return optionFunc(func(o *options) {
		o.vaultSubPath = vaultSubPath
	})
}

// Store is the store for the wallet held encrypted on Amazon S3.
type Store struct {
	client       *api.Client
	jwt          string
	passphrase   []byte
	role         string
	vaultSubPath string
}

// New creates a new Vault backed store.
// This takes the following options:
//  - region: a string specifying the Amazon S3 region, defaults to "us-east-1", set with WithRegion()
//  - id: a byte array specifying an identifying key for the store, defaults to nil, set with WithID()
// This expects the access credentials to be in a standard place, e.g. ~/.aws/credentials
func New(opts ...Option) (wtypes.Store, error) {
	options := options{
		vaultAddress: "http://vault.vault:8200",
		role:         "eth",
		vaultSubPath: "eth",
	}
	for _, o := range opts {
		o.apply(&options)
	}

	client, err := api.NewClient(&api.Config{
		Address: options.vaultAddress,
	})

	if err != nil {
		//log.Printf("error creating new client %v", err)
		return nil, err
	}

	jwt, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")

	if err != nil {
		//log.Printf("error reading serviceaccount token %v", err)
		return nil, err
	}

	return &Store{
		client:       client,
		jwt:          string(jwt),
		passphrase:   options.passphrase,
		role:         options.role,
		vaultSubPath: options.vaultSubPath,
	}, nil
}

func (s *Store) Authorize() error {
	client := s.client

	config := map[string]interface{}{
		"role": s.role,
		// Have to convert this into a string to compact the jwt
		"jwt": s.jwt,
	}

	//log.Printf("attempting to write with role: %v and jtw: %v", s.role, s.jwt)

	resp, err := client.Logical().Write("auth/kubernetes/login", config)

	if err != nil {
		//log.Printf("error writing config to auth/kubernetes/login: %v", err)
		return err
	}

	client.SetToken(resp.Auth.ClientToken)
	//log.Printf("headers: %v", client.Headers())
	//log.Printf("set token as %v", resp.Auth.ClientToken)

	return nil
}

// Name returns the name of this store.
func (s *Store) Name() string {
	return "vault"
}

// Location returns the location of this store.
func (s *Store) Location() string {
	return s.vaultSubPath
}
