package tlogm

import (
	"sync"

	"github.com/honeycombio/libhoney-go/transmission"
	"github.com/ibrt/golang-utils/memz"
)

var (
	_ transmission.Sender = (*MockSender)(nil)
)

// MockSender implements the transmission.Sender interface.
type MockSender struct {
	m      *sync.Mutex
	events []*transmission.Event
	c      chan transmission.Response
}

// NewMockSender initializes a new MockSender.
func NewMockSender() *MockSender {
	return &MockSender{
		m:      &sync.Mutex{},
		events: make([]*transmission.Event, 0),
		c:      make(chan transmission.Response),
	}
}

// GetEvents returns the events.
func (s *MockSender) GetEvents() []*transmission.Event {
	s.m.Lock()
	defer s.m.Unlock()
	return memz.ShallowCopySlice(s.events)
}

// ClearEvents clears the events.
func (s *MockSender) ClearEvents() {
	s.m.Lock()
	defer s.m.Unlock()
	s.events = make([]*transmission.Event, 0)
}

// Add implements the transmission.Sender interface.
func (s *MockSender) Add(e *transmission.Event) {
	s.m.Lock()
	defer s.m.Unlock()
	s.events = append(s.events, e)
}

// Start implements the transmission.Sender interface.
func (s *MockSender) Start() error {
	return nil
}

// Stop implements the transmission.Sender interface.
func (s *MockSender) Stop() error {
	s.m.Lock()
	defer s.m.Unlock()

	s.events = nil
	close(s.c)
	return nil
}

// Flush implements the transmission.Sender interface.
func (s *MockSender) Flush() error {
	return nil
}

// TxResponses implements the transmission.Sender interface.
func (s *MockSender) TxResponses() chan transmission.Response {
	return s.c
}

// SendResponse implements the transmission.Sender interface.
func (s *MockSender) SendResponse(response transmission.Response) bool {
	s.m.Lock()
	defer s.m.Unlock()

	select {
	case s.c <- response:
		return false
	default:
		return true
	}
}
