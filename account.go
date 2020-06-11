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
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// StoreAccount stores an account.  It will fail if it cannot store the data.
// Note this will overwrite an existing account with the same ID.  It will not, however, allow multiple accounts with the same
// name to co-exist in the same wallet.
func (s *Store) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
	client := s.client
	// Ensure the wallet exists
	_, err := s.RetrieveWalletByID(walletID)
	if err != nil {
		return errors.New("unknown wallet")
	}

	// See if an account with this name already exists
	existingAccount, err := s.RetrieveAccount(walletID, accountID)
	if err == nil {
		// It does; they need to have the same ID for us to overwrite it
		info := &struct {
			ID string `json:"uuid"`
		}{}
		err := json.Unmarshal(existingAccount, info)
		if err != nil {
			return err
		}
		if info.ID != accountID.String() {
			return errors.New("account already exists")
		}
	}

	path := s.accountPath(walletID.String(), accountID.String())

	_, err = client.Logical().WriteBytes(path, data)

	if err != nil {
		return errors.Wrap(err, "failed to store key")
	}
	return nil
}

// RetrieveAccount retrieves account-level data.  It will fail if it cannot retrieve the data.
func (s *Store) RetrieveAccount(walletID uuid.UUID, accountID uuid.UUID) ([]byte, error) {
	client := s.client
	path := s.accountPath(walletID.String(), accountID.String())

	secret, err := client.Logical().Read(path)

	if err != nil {
		return nil, err
	}

	byteData, err := json.Marshal(secret.Data)

	if err != nil {
		return nil, err
	}

	return byteData, nil
}

// RetrieveAccounts retrieves all account-level data for a wallet.
func (s *Store) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
	client := s.client
	path := s.walletPath(walletID.String())
	ch := make(chan []byte, 1024)
	go func() {
		secret, err := client.Logical().List(path)

		if err != nil {
			return
		}

		// Discard this error for now
		// TODO: Do something with the error
		accounts, _ := secret.Data["keys"].([]interface{})

		for _, account := range accounts {
			if account.(string) != "index" && account.(string) != walletID.String() {

				// Quietly skip these errors
				// TODO: Handle errors better through the channel
				secret, err := client.Logical().Read(s.accountPath(walletID.String(), account.(string)))

				if err != nil {
					continue
				}

				byteData, err := json.Marshal(secret.Data)

				if err != nil {
					continue
				}

				data, err := s.decryptIfRequired(byteData)

				if err != nil {
					continue
				}
				ch <- data
			}
		}
		close(ch)
	}()
	return ch
}
