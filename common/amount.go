// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package common

import (
	"fmt"
)

// Amount describes a method of converting an Amount to something
// other than the base unit of an asimov.
type Amount int64

// String returns the unit as a string.  For recognized units, the SI
// prefix is used, or "Xing" for the base unit.
func (u Amount) String() string {
	if u > 1000 {
		return fmt.Sprintf("%s,%03d", (u/Amount(1000)).String(), u % 1000)
	}
	return fmt.Sprintf("%d", u % 1000)
}
