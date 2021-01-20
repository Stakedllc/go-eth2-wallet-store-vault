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
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// StoreAccountsIndex stores the account index.
func (s *Store) StoreAccountsIndex(walletID uuid.UUID, data []byte) error {
	var err error

	// Do not encrypt empty index.
	if len(data) != 2 {
		data, err = s.encryptIfRequired(data)
		if err != nil {
			return err
		}
	}

	path := filepath.Join(s.localPath, "index.json")
	err = ioutil.WriteFile(path, data, 0600)

	if err != nil {
		return errors.Wrap(err, "failed to store key")
	}
	return nil
}

// RetrieveAccountsIndex retrieves the account index.
func (s *Store) RetrieveAccountsIndex(walletID uuid.UUID) ([]byte, error) {
	path := filepath.Join(s.localPath, "index.json")
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	// Do not decrypt empty index.
	if len(data) == 2 {
		return data, nil
	}

	if data, err = s.decryptIfRequired(data); err != nil {
		return nil, err
	}
	return data, nil
}
