package user_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type MockClientConn struct {
	mockMethods map[string]any
	status      string
}

func (m *MockClientConn) Invoke(
	ctx context.Context,
	method string,
	args any,
	reply any,
	opts ...grpc.CallOption,
) error {
	if m.mockMethods == nil {
		m.mockMethods = make(map[string]any)
	}
	if strings.Contains(method, "MakeEvent") || strings.Contains(method, "DeleteEvent") {
		m.mockMethods[method] = func() any {
			return &eventctrl.EventIdAvail{EventId: 12345}
		}
	}
	if strings.Contains(method, "GetEvent") {
		m.mockMethods[method] = func() any {
			return &eventctrl.Event{
				SenderId: 12345,
				EventId:  12345,
				Time:     12345,
				Name:     "test",
			}
		}
	}
	if methodFunc, ok := m.mockMethods[method]; ok {
		if fn, ok := methodFunc.(func() any); ok {
			result := fn()
			m.status = fmt.Sprintln(result)

			switch v := result.(type) {
			case *eventctrl.EventIdAvail:
				*reply.(*eventctrl.EventIdAvail) = *v
			case *eventctrl.Event:
				*reply.(*eventctrl.Event) = *v
			default:
				return errors.New("unexpected result type")
			}
		} else {
			return errors.New("unexpected function type")
		}
	} else {
		return errors.New("method not found")
	}
	return nil
}

func (m *MockClientConn) NewStream(
	ctx context.Context,
	desc *grpc.StreamDesc,
	method string,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	return nil, nil
}

func TestMakeEvent(t *testing.T) {
	mockConn := &MockClientConn{}
	client := eventctrl.NewEventsClient(mockConn)

	ctx := context.Background()
	request := &eventctrl.MakeEventRequest{
		// TODO: fill
		// struct
		SenderId: 1,
		Time:     1432423,
		Name:     "User1",
	}

	response, err := client.MakeEvent(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(12345), response.EventId)
}

func TestGetEvent(t *testing.T) {
	mockConn := &MockClientConn{}
	client := eventctrl.NewEventsClient(mockConn)

	ctx := context.Background()
	request := &eventctrl.GetEventRequest{
		// TODO: fill
		// struct
		SenderId: 1,
		EventId:  1,
	}

	response, err := client.GetEvent(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(12345), response.SenderId)
	assert.Equal(t, int64(12345), response.EventId)
	assert.Equal(t, int64(12345), response.Time)
	assert.Equal(t, "test", response.Name)
}
