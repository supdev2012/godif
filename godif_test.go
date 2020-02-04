/*
 * Copyright (c) 2018-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package godif

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncErrorOnNoImplementation(t *testing.T) {
	Reset()
	var injectedFunc func(x int, y int) int

	Require(&injectedFunc)
	errs := ResolveAll()

	if _, ok := errs[0].(*EImplementationNotProvided); ok && len(errs) == 1 {
		fmt.Println(errs)
	} else {
		t.Fatal(errs)
	}

	assert.Nil(t, injectedFunc)
}
