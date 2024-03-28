package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/Raimguzhinov/simple-grpc/internal/models"
	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type testStruct struct {
	name     string
	actual   models.Event
	expected models.Event
}

type MockClientConn struct {
	mockMethods map[string]func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error
}

func NewMockClientConn() *MockClientConn {
	return &MockClientConn{
		make(map[string]func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error),
	}
}

func (m *MockClientConn) MockMethod(
	method string,
	mockFunc func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error,
) {
	m.mockMethods[method] = mockFunc
}

func (m *MockClientConn) Invoke(
	ctx context.Context,
	method string,
	args any,
	reply any,
	opts ...grpc.CallOption,
) error {
	if methodFunc, ok := m.mockMethods[method]; ok {
		err := methodFunc(ctx, args, reply)
		if err != nil {
			return err
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
	mockConn := NewMockClientConn()
	currentTime := time.Now().UnixMilli()

	client := eventctrl.NewEventsClient(mockConn)

	mockConn.MockMethod(
		eventctrl.Events_MakeEvent_FullMethodName,
		func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error {
			req := args.(*eventctrl.MakeEventRequest)
			resp := reply.(*eventctrl.EventIdAvail)
			resp.EventId = 12345
			_ = req
			return nil
		},
	)

	ctx := context.Background()

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name: "Test 1",
			actual: models.Event{
				SenderID: 1,
				Time:     currentTime,
				Name:     "test",
			},
			expected: models.Event{ID: 12345},
		},
		testStruct{
			name: "Test 2",
			actual: models.Event{
				SenderID: 2,
				Time:     currentTime,
				Name:     "test",
			},
			expected: models.Event{ID: 12345},
		},
	)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			request := &eventctrl.MakeEventRequest{
				SenderId: testCase.actual.SenderID,
				Time:     testCase.actual.Time,
				Name:     testCase.actual.Name,
			}
			response, err := client.MakeEvent(ctx, request)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.expected.ID, response.EventId)
		})
	}
}

func TestGetEvent(t *testing.T) {
	mockConn := NewMockClientConn()
	currentTime := time.Now().UnixMilli()

	client := eventctrl.NewEventsClient(mockConn)

	mockConn.MockMethod(
		eventctrl.Events_GetEvent_FullMethodName,
		func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error {
			req := args.(*eventctrl.GetEventRequest)
			resp := reply.(*eventctrl.Event)

			resp.EventId = req.EventId
			resp.SenderId = req.SenderId
			resp.Time = currentTime
			resp.Name = "test"
			return nil
		},
	)

	ctx := context.Background()

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name: "Test 1",
			actual: models.Event{
				SenderID: 1,
				ID:       1,
			},
			expected: models.Event{
				SenderID: 1,
				ID:       1,
				Time:     currentTime,
				Name:     "test",
			},
		},
		testStruct{
			name: "Test 2",
			actual: models.Event{
				SenderID: 2,
				ID:       2,
			},
			expected: models.Event{
				SenderID: 2,
				ID:       2,
				Time:     currentTime,
				Name:     "test",
			},
		},
	)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			request := &eventctrl.GetEventRequest{
				SenderId: testCase.actual.SenderID,
				EventId:  testCase.actual.ID,
			}
			response, err := client.GetEvent(ctx, request)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.expected.SenderID, response.SenderId)
			assert.Equal(t, testCase.expected.ID, response.EventId)
			assert.Equal(t, testCase.expected.Time, response.Time)
			assert.Equal(t, testCase.expected.Name, response.Name)
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	mockConn := NewMockClientConn()
	client := eventctrl.NewEventsClient(mockConn)

	mockConn.MockMethod(
		eventctrl.Events_DeleteEvent_FullMethodName,
		func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error {
			req := args.(*eventctrl.DeleteEventRequest)
			resp := reply.(*eventctrl.EventIdAvail)
			resp.EventId = req.EventId
			return nil
		},
	)

	ctx := context.Background()

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name: "Test 1",
			actual: models.Event{
				SenderID: 1,
				ID:       1,
			},
			expected: models.Event{ID: 1},
		},
		testStruct{
			name: "Test 2",
			actual: models.Event{
				SenderID: 2,
				ID:       2,
			},
			expected: models.Event{ID: 2},
		},
	)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			request := &eventctrl.DeleteEventRequest{
				SenderId: testCase.actual.SenderID,
				EventId:  testCase.actual.ID,
			}
			response, err := client.DeleteEvent(ctx, request)
			if err != nil {
				t.Fatal(err)
			}
			assert.EqualValues(t, testCase.expected.ID, response.EventId)
		})
	}
}
