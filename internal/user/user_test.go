package user_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Raimguzhinov/simple-grpc/internal/user"
	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type testStruct struct {
	name     string
	senderId int64
	promt    string
	expected string
}

type MockClientConn struct {
	mockMethods map[string]func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error
	mockStreams map[string]func(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error)
}

func NewMockClientConn() *MockClientConn {
	return &MockClientConn{
		mockMethods: make(
			map[string]func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error,
		),
		mockStreams: make(
			map[string]func(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error),
		),
	}
}

func (m *MockClientConn) MockMethod(
	method string,
	mockFunc func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error,
) {
	m.mockMethods[method] = mockFunc
}

func (m *MockClientConn) MockStream(
	method string,
	mockFunc func(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error),
) {
	m.mockStreams[method] = mockFunc
}

func (m *MockClientConn) Invoke(
	ctx context.Context,
	method string,
	args any,
	reply any,
	opts ...grpc.CallOption,
) error {
	if methodFunc, ok := m.mockMethods[method]; ok {
		err := methodFunc(ctx, args, reply, opts...)
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
	if streamFunc, ok := m.mockStreams[method]; ok {
		return streamFunc(ctx, desc, method, opts...)
	} else {
		return nil, errors.New("stream method not found")
	}
}

func TestMakeEvent(t *testing.T) {
	mockConn := NewMockClientConn()
	currentTime := time.Now()
	customID := uuid.New()
	binCustomID, _ := customID.MarshalBinary()

	client := eventctrl.NewEventsClient(mockConn)
	usr := user.NewUser(&client)

	mockConn.MockMethod(
		eventctrl.Events_MakeEvent_FullMethodName,
		func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error {
			req := args.(*eventctrl.MakeEventRequest)
			resp := reply.(*eventctrl.EventIdAvail)
			resp.EventId = binCustomID
			_ = req
			return nil
		},
	)

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name:     "Test 1",
			senderId: 1,
			promt:    fmt.Sprintf("%s %s", currentTime.Format(time.DateTime), "User1"),
			expected: fmt.Sprintf("Created {%s}\n", customID.String()),
		},
		testStruct{
			name:     "Test 2",
			senderId: 2,
			promt:    fmt.Sprintf("%s %s", currentTime.Format(time.DateTime), "User2"),
			expected: fmt.Sprintf("Created {%s}\n", customID.String()),
		},
		testStruct{
			name:     "Test 3",
			senderId: 3,
			promt:    fmt.Sprintf("%s %s", "12/09-2003 2:00", "User3"),
			expected: fmt.Sprintf("%s\n", user.ErrBadDateTime),
		},
		testStruct{
			name:     "Test 4",
			senderId: 4,
			promt:    fmt.Sprintf("%s.", "User4"),
			expected: fmt.Sprintf("%s\n", user.ErrBadFormat),
		},
	)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			usr.EventMaker(testCase.senderId, testCase.promt, buf)
			require.Equal(t, testCase.expected, buf.String())
		})
	}
}

func TestGetEvent(t *testing.T) {
	mockConn := NewMockClientConn()
	currentTime := time.Now()
	customID := uuid.New()

	client := eventctrl.NewEventsClient(mockConn)
	usr := user.NewUser(&client)

	mockConn.MockMethod(
		eventctrl.Events_GetEvent_FullMethodName,
		func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error {
			req := args.(*eventctrl.GetEventRequest)
			resp := reply.(*eventctrl.Event)

			resp.EventId = req.EventId
			resp.SenderId = req.SenderId
			resp.Time = currentTime.UnixMilli()
			resp.Name = "test"
			return nil
		},
	)

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name:     "Test 1",
			senderId: 1,
			promt:    customID.String(),
			expected: fmt.Sprintf(
				"Event {\n  senderId: %d\n  eventId: %s\n  time: %s\n  name: '%s'\n}\n",
				1,
				customID.String(),
				currentTime.Format(time.DateTime),
				"test",
			),
		},
		testStruct{
			name:     "Test 2",
			senderId: 2,
			promt:    fmt.Sprintf("%s-0*bb0", customID.String()),
			expected: fmt.Sprintf("%s\n", user.ErrBadFormat),
		},
	)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			usr.EventGetter(testCase.senderId, testCase.promt, buf)
			require.Equal(t, testCase.expected, buf.String())
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	mockConn := NewMockClientConn()
	customID := uuid.New()

	client := eventctrl.NewEventsClient(mockConn)
	usr := user.NewUser(&client)

	mockConn.MockMethod(
		eventctrl.Events_DeleteEvent_FullMethodName,
		func(ctx context.Context, args, reply any, opts ...grpc.CallOption) error {
			req := args.(*eventctrl.DeleteEventRequest)
			resp := reply.(*eventctrl.EventIdAvail)
			resp.EventId = req.EventId
			return nil
		},
	)

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name:     "Test 1",
			senderId: 1,
			promt:    customID.String(),
			expected: fmt.Sprintf("Deleted {%s}\n", customID.String()),
		},
		testStruct{
			name:     "Test 2",
			senderId: 2,
			promt:    fmt.Sprintf("%s-0*bb0", customID.String()),
			expected: fmt.Sprintf("%s\n", user.ErrBadFormat),
		},
	)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			usr.EventDeleter(testCase.senderId, testCase.promt, buf)
			require.Equal(t, testCase.expected, buf.String())
		})
	}
}

func TestGetEvents(t *testing.T) {
	mockConn := NewMockClientConn()
	customID := uuid.New()
	currentTime := time.Now()
	mockStream := NewMockClientStream(customID, currentTime)

	client := eventctrl.NewEventsClient(mockConn)
	usr := user.NewUser(&client)

	mockConn.MockStream(
		eventctrl.Events_GetEvents_FullMethodName,
		func(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return mockStream, nil
		},
	)

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name:     "Test 1",
			senderId: 1,
			promt: fmt.Sprintf(
				"%s %s",
				currentTime.AddDate(0, -1, 0).Format(time.DateTime),
				currentTime.AddDate(0, 1, 0).Format(time.DateTime),
			),
			expected: fmt.Sprintf(
				"Event {\n  senderId: %d\n  eventId: %s\n  time: %s\n  name: '%s'\n}\nEvent {\n  senderId: %d\n  eventId: %s\n  time: %s\n  name: '%s'\n}\n",
				1,
				customID.String(),
				currentTime.Format(time.DateTime),
				"test",
				1,
				customID.String(),
				currentTime.Format(time.DateTime),
				"test",
			),
		},
	)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			usr.EventsGetter(testCase.senderId, testCase.promt, buf)
			require.Equal(t, testCase.expected, buf.String())
		})
	}
}

type MockClientStream interface {
	Header() (metadata.MD, error)
	Trailer() metadata.MD
	CloseSend() error
	Context() context.Context
	SendMsg(m any) error
	RecvMsg(m any) error
}

type mockClientStream struct {
	grpc.ClientStream
	resvCount   int
	respStr     string
	customID    []byte
	currentTime time.Time
}

func NewMockClientStream(eventID uuid.UUID, currentTime time.Time) MockClientStream {
	binUID, _ := eventID.MarshalBinary()
	return &mockClientStream{
		customID:    binUID,
		currentTime: currentTime,
	}
}

func (m *mockClientStream) SendMsg(msg any) error {
	m.respStr = fmt.Sprintf("%+v\n", msg)
	return nil
}

func (m *mockClientStream) CloseSend() error {
	return nil
}

func (m *mockClientStream) RecvMsg(msg any) error {
	keys := strings.Split(m.respStr, ":")
	valueSender := strings.Split(keys[1], " ")
	senderID, _ := strconv.Atoi(valueSender[0])

	m.resvCount++
	msg.(*eventctrl.Event).Name = "test"
	msg.(*eventctrl.Event).EventId = m.customID
	msg.(*eventctrl.Event).Time = m.currentTime.UnixMilli()
	msg.(*eventctrl.Event).SenderId = int64(senderID)
	if m.resvCount >= 3 {
		return io.EOF
	}
	return nil
}
