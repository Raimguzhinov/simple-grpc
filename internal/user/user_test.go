package user_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type MockClientConn struct{}

func (m *MockClientConn) Invoke(
	ctx context.Context,
	method string,
	args any,
	reply any,
	opts ...grpc.CallOption,
) error {
	// TODO: add logic
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
	}

	response, err := client.MakeEvent(ctx, request)
	// TODO:
	// check and assert
	_ = response
	_ = err
}
