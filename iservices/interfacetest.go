/*
 * Copyright (c) 2019-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package iservices

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untillpro/godif"
)

var declareImplementation func()

// TestAll s.e.
func TestAll(t *testing.T, declare func()) {
	declareImplementation = declare
	t.Run("TestBasicUsage", TestBasicUsage)
	t.Run("TestFailedInit", TestFailedInit)
	t.Run("TestFailedStart", TestFailedStart)
}

// TestBasicUsage s.e.
func TestBasicUsage(t *testing.T) {

	// We will need InitAndStart & StopAndFinit functions

	godif.Require(&InitAndStart)
	godif.Require(&StopAndFinit)

	// Declare passed implementation

	declareImplementation()

	// Provide own services

	s1 := &myService{Name: "Service1"}
	s2 := &myService{Name: "Service2"}
	godif.ProvideSliceElement(&Services, s1)
	godif.ProvideSliceElement(&Services, s2)

	// Resolve all

	errs := godif.ResolveAll()
	defer godif.Reset()
	assert.Nil(t, errs)
	fmt.Println("errs=", errs)

	// Init and start services

	ctx := context.Background()
	ctx, err := InitAndStart(ctx)
	assert.Nil(t, err)
	//2 means "started"
	assert.Equal(t, 2, s1.State)
	assert.Equal(t, 2, s2.State)

	// Context should be proper initialized

	fmt.Println("Ctx=", ctx)

	//Make sure that value provided by service exist in ctx
	assert.True(t, ctx.Value(ctxKeyType("Service1")).(bool))
	assert.True(t, s1.ctxValue.(bool))

	assert.True(t, ctx.Value(ctxKeyType("Service2")).(bool))
	assert.True(t, s2.ctxValue.(bool))

	assert.Nil(t, ctx.Value(ctxKeyType("Service3")))

	// Stop and finit services

	StopAndFinit(ctx)
	//State must be 0
	assert.Equal(t, 0, s1.State)
	assert.Equal(t, 0, s2.State)

}

// TestFailedInit s.e.
func TestFailedInit(t *testing.T) {

	// Declare iservices requirements and implementation

	godif.Require(&InitAndStart)
	godif.Require(&StopAndFinit)
	declareImplementation()

	// Provide services, s2 will fail on start
	s1 := &myService{Name: "Service1"}
	s2 := &myService{Name: "Service2", Failinit: true}
	godif.ProvideSliceElement(&Services, s1)
	godif.ProvideSliceElement(&Services, s2)
	errs := godif.ResolveAll()
	defer godif.Reset()
	assert.Nil(t, errs)

	// Init and start services

	ctx := context.Background()
	ctx, err := InitAndStart(ctx)
	assert.NotNil(t, err)
	fmt.Println("err=", err)
	assert.True(t, strings.Contains(err.Error(), "Service2"))
	assert.False(t, strings.Contains(err.Error(), "Service1"))
	assert.Equal(t, 1, s1.State)
	assert.Equal(t, 0, s2.State)

	// Stop and finit services

	StopAndFinit(ctx)
	assert.Equal(t, 0, s1.State)
	assert.Equal(t, 0, s2.State)

}

// TestFailedStart s.e.
func TestFailedStart(t *testing.T) {

	// Declare iservices requirements and implementation

	godif.Require(&InitAndStart)
	godif.Require(&StopAndFinit)
	declareImplementation()

	// Provide services, s2 will fail on start
	s1 := &myService{Name: "Service1"}
	s2 := &myService{Name: "Service2", Failstart: true}
	godif.ProvideSliceElement(&Services, s1)
	godif.ProvideSliceElement(&Services, s2)
	errs := godif.ResolveAll()
	defer godif.Reset()
	assert.Nil(t, errs)

	// Init and start services

	ctx := context.Background()
	ctx, err := InitAndStart(ctx)
	assert.NotNil(t, err)
	fmt.Println("err=", err)
	assert.True(t, strings.Contains(err.Error(), "Service2"))
	assert.False(t, strings.Contains(err.Error(), "Service1"))
	assert.Equal(t, 2, s1.State)
	assert.Equal(t, 1, s2.State)

	// Stop and finit services

	StopAndFinit(ctx)
	assert.Equal(t, 0, s1.State)
	assert.Equal(t, 0, s2.State)

}

type myService struct {
	Name      string
	State     int // 0, 1(inited), 2(started), 3 (stopped)
	Failstart bool
	Failinit  bool
	ctxValue  interface{}
}

type ctxKeyType string

func (s *myService) Init(ctx context.Context) (context.Context, error) {
	if s.Failinit {
		fmt.Println(s.Name, "Init fails")
		return ctx, errors.New(s.Name + ":" + "Init fails")
	}
	s.State++
	fmt.Println(s.Name, "Inited")
	ctx = context.WithValue(ctx, ctxKeyType(s.Name), true)
	return ctx, nil
}

func (s *myService) Start(ctx context.Context) error {
	if s.Failstart {
		fmt.Println(s.Name, "Start fails")
		return errors.New(s.Name + ":" + "Start fails")
	}

	s.ctxValue = ctx.Value(ctxKeyType(s.Name))

	s.State++
	fmt.Println(s.Name, "Started")
	return nil
}

func (s *myService) Stop(ctx context.Context) {
	s.State--
	fmt.Println(s.Name, "Stopped")
}

func (s *myService) Finit(ctx context.Context) {
	s.State--
	fmt.Println(s.Name, "Finited")
	s.Failinit = false
	s.Failstart = false
}

func (s *myService) String() string {
	return "I'm service " + s.Name
}