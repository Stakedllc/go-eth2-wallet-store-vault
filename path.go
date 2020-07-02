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

package vault

import (
	"fmt"
)

func (s *Store) walletsPath() string {
	return fmt.Sprintf("/secret/%s", s.Location())
}

func (s *Store) walletPath(walletID string) string {
	return fmt.Sprintf("/secret/%s/%s", s.Location(), walletID)
}

func (s *Store) walletHeaderPath(walletID string) string {
	return fmt.Sprintf("/secret/%s/%s/%s", s.Location(), walletID, walletID)
}

func (s *Store) accountPath(walletID string, accountID string) string {
	return fmt.Sprintf("/secret/%s/%s/%s", s.Location(), walletID, accountID)
}

func (s *Store) walletIndexPath(walletID string) string {
	return fmt.Sprintf("/secret/%s/%s/index", s.Location(), walletID)
}
