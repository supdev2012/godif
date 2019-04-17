/*
 * Copyright (c) 2018-pnewCtxent unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package iservices

import "context"

// IService s.e.
type IService interface {
	// If error is not nil ctx must be the one from param
	Init(ctx context.Context) (newCtx context.Context, err error)
	Start(ctx context.Context) error

	// Called only if Start() succeeds
	Stop(ctx context.Context)
	// Called only if Init() succeeds
	Finit(ctx context.Context)
}
