/*
 * Copyright (c) 2018-pnewCtxent unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package iservices

import (
	"context"
	"sync"
)


var inited = []IService{}
var started = []IService{}

func initAndStartImpl(ctx context.Context) (newCtx context.Context, err error) {
	newCtx = ctx
	for _, service := range Services {
		newCtx, err = service.Init(newCtx)
		if nil != err {
			return newCtx, err
		}
		inited = append(inited, service)
	}

	for _, service := range Services {
		err := service.Start(newCtx)
		if nil != err {
			return newCtx, err
		}
		started = append(started, service)
	}
	return newCtx, nil
}

func stopAndFinitImpl(ctx context.Context) {

	var wg sync.WaitGroup

	for _, service := range started {
		s := service
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Stop(ctx)
		}()
	}
	wg.Wait()
	started = []IService{}

	for i := len(inited) - 1; i >= 0; i-- {
		inited[i].Finit(ctx)
	}
	inited = []IService{}
}