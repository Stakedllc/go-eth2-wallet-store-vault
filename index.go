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
	b64 "encoding/base64"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// StoreAccountsIndex stores the account index.
func (s *Store) StoreAccountsIndex(walletID uuid.UUID, data []byte) error {
	s.Authorize()

	client := s.client
	var err error
	var structuredData map[string]interface{}

	// Do not encrypt empty index.
	if len(data) != 2 {

		data, err = s.encryptIfRequired(data)
		if err != nil {
			return err
		}

		structuredData = map[string]interface{}{
			"data": data,
		}
	} else {
		structuredData = map[string]interface{}{
			"data": data,
		}
	}

	path := s.walletIndexPath(walletID.String())

	_, err = client.Logical().Write(path, structuredData)

	if err != nil {
		return errors.Wrap(err, "failed to store key")
	}
	return nil
}

// RetrieveAccountsIndex retrieves the account index.
func (s *Store) RetrieveAccountsIndex(walletID uuid.UUID) ([]byte, error) {
	s.Authorize()

	client := s.client
	path := s.walletIndexPath(walletID.String())

	secret, err := client.Logical().Read(path)

	if err != nil {
		return nil, err
	}

	byteData, _ := b64.StdEncoding.DecodeString(secret.Data["data"].(string))

	data, err := s.decryptIfRequired(byteData)
	if err != nil {
		return nil, err
	}

	return data, nil
}
