// This file was generated by github.com/nelsam/hel.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package octo_test

import (
	"context"
	"time"

	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type mockOctoService struct {
	FetchAllOctosCalled chan bool
	FetchAllOctosInput  struct {
		Ctx chan context.Context
	}
	FetchAllOctosOutput struct {
		Octos chan []data.Octo
		Err   chan error
	}
	FetchOctoByNameCalled chan bool
	FetchOctoByNameInput  struct {
		Ctx  chan context.Context
		Name chan string
	}
	FetchOctoByNameOutput struct {
		Octo chan data.Octo
		Err  chan error
	}
	CreateOctoCalled chan bool
	CreateOctoInput  struct {
		Ctx    chan context.Context
		OctoIn chan data.Octo
	}
	CreateOctoOutput struct {
		OctoOut chan data.Octo
		Err     chan error
	}
	DeleteOctoByNameCalled chan bool
	DeleteOctoByNameInput  struct {
		Ctx  chan context.Context
		Name chan string
	}
	DeleteOctoByNameOutput struct {
		Err chan error
	}
}

func newMockOctoService() *mockOctoService {
	m := &mockOctoService{}
	m.FetchAllOctosCalled = make(chan bool, 100)
	m.FetchAllOctosInput.Ctx = make(chan context.Context, 100)
	m.FetchAllOctosOutput.Octos = make(chan []data.Octo, 100)
	m.FetchAllOctosOutput.Err = make(chan error, 100)
	m.FetchOctoByNameCalled = make(chan bool, 100)
	m.FetchOctoByNameInput.Ctx = make(chan context.Context, 100)
	m.FetchOctoByNameInput.Name = make(chan string, 100)
	m.FetchOctoByNameOutput.Octo = make(chan data.Octo, 100)
	m.FetchOctoByNameOutput.Err = make(chan error, 100)
	m.CreateOctoCalled = make(chan bool, 100)
	m.CreateOctoInput.Ctx = make(chan context.Context, 100)
	m.CreateOctoInput.OctoIn = make(chan data.Octo, 100)
	m.CreateOctoOutput.OctoOut = make(chan data.Octo, 100)
	m.CreateOctoOutput.Err = make(chan error, 100)
	m.DeleteOctoByNameCalled = make(chan bool, 100)
	m.DeleteOctoByNameInput.Ctx = make(chan context.Context, 100)
	m.DeleteOctoByNameInput.Name = make(chan string, 100)
	m.DeleteOctoByNameOutput.Err = make(chan error, 100)
	return m
}
func (m *mockOctoService) FetchAllOctos(ctx context.Context) (octos []data.Octo, err error) {
	m.FetchAllOctosCalled <- true
	m.FetchAllOctosInput.Ctx <- ctx
	return <-m.FetchAllOctosOutput.Octos, <-m.FetchAllOctosOutput.Err
}
func (m *mockOctoService) FetchOctoByName(ctx context.Context, name string) (octo data.Octo, err error) {
	m.FetchOctoByNameCalled <- true
	m.FetchOctoByNameInput.Ctx <- ctx
	m.FetchOctoByNameInput.Name <- name
	return <-m.FetchOctoByNameOutput.Octo, <-m.FetchOctoByNameOutput.Err
}
func (m *mockOctoService) CreateOcto(ctx context.Context, octoIn data.Octo) (octoOut data.Octo, err error) {
	m.CreateOctoCalled <- true
	m.CreateOctoInput.Ctx <- ctx
	m.CreateOctoInput.OctoIn <- octoIn
	return <-m.CreateOctoOutput.OctoOut, <-m.CreateOctoOutput.Err
}
func (m *mockOctoService) DeleteOctoByName(ctx context.Context, name string) (err error) {
	m.DeleteOctoByNameCalled <- true
	m.DeleteOctoByNameInput.Ctx <- ctx
	m.DeleteOctoByNameInput.Name <- name
	return <-m.DeleteOctoByNameOutput.Err
}

type mockContext struct {
	DeadlineCalled chan bool
	DeadlineOutput struct {
		Deadline chan time.Time
		Ok       chan bool
	}
	DoneCalled chan bool
	DoneOutput struct {
		Ret0 chan (<-chan struct{})
	}
	ErrCalled chan bool
	ErrOutput struct {
		Ret0 chan error
	}
	ValueCalled chan bool
	ValueInput  struct {
		Key chan interface{}
	}
	ValueOutput struct {
		Ret0 chan interface{}
	}
}

func newMockContext() *mockContext {
	m := &mockContext{}
	m.DeadlineCalled = make(chan bool, 100)
	m.DeadlineOutput.Deadline = make(chan time.Time, 100)
	m.DeadlineOutput.Ok = make(chan bool, 100)
	m.DoneCalled = make(chan bool, 100)
	m.DoneOutput.Ret0 = make(chan (<-chan struct{}), 100)
	m.ErrCalled = make(chan bool, 100)
	m.ErrOutput.Ret0 = make(chan error, 100)
	m.ValueCalled = make(chan bool, 100)
	m.ValueInput.Key = make(chan interface{}, 100)
	m.ValueOutput.Ret0 = make(chan interface{}, 100)
	return m
}
func (m *mockContext) Deadline() (deadline time.Time, ok bool) {
	m.DeadlineCalled <- true
	return <-m.DeadlineOutput.Deadline, <-m.DeadlineOutput.Ok
}
func (m *mockContext) Done() <-chan struct{} {
	m.DoneCalled <- true
	return <-m.DoneOutput.Ret0
}
func (m *mockContext) Err() error {
	m.ErrCalled <- true
	return <-m.ErrOutput.Ret0
}
func (m *mockContext) Value(key interface{}) interface{} {
	m.ValueCalled <- true
	m.ValueInput.Key <- key
	return <-m.ValueOutput.Ret0
}