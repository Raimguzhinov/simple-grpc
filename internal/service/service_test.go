package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Raimguzhinov/simple-grpc/internal/models"
	"github.com/Raimguzhinov/simple-grpc/internal/service"
	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

func TestMakeEvent(t *testing.T) {
	// Arrange:
	s := service.RunEventsService()
	currentTime := time.Now().UnixMilli()
	// nowTime := time.Now()
	// nowTime2 := time.UnixMilli(currentTime).UTC()
	// require.Equal(t, currentTime, nowTime.UnixMilli())
	// require.Equal(t, nowTime, nowTime2)

	testTable := []struct {
		actual   models.Events
		expected models.Events
	}{
		{
			actual: models.Events{
				SenderID: 1,
				Time:     currentTime,
				Name:     "User1",
			},
			expected: models.Events{
				ID: 1,
			},
		},
		{
			actual: models.Events{
				SenderID: 1,
				Time:     currentTime,
				Name:     "User1Again",
			},
			expected: models.Events{
				ID: 2,
			},
		},
		{
			actual: models.Events{
				SenderID: 2,
				Time:     currentTime,
				Name:     "User2",
			},
			expected: models.Events{
				ID: 1,
			},
		},
		{
			actual: models.Events{
				SenderID: 1,
				Time:     currentTime,
				Name:     "User1AgainAgain",
			},
			expected: models.Events{
				ID: 3,
			},
		},
		{
			actual: models.Events{
				SenderID: 2,
				Time:     currentTime,
				Name:     "TestUser2Again",
			},
			expected: models.Events{
				ID: 2,
			},
		},
	}

	// Act:
	for _, testCase := range testTable {
		req := &eventmanager.MakeEventRequest{
			SenderId: testCase.actual.SenderID,
			Time:     testCase.actual.Time,
			Name:     testCase.actual.Name,
		}
		resp, err := s.MakeEvent(context.Background(), req)
		t.Logf("\nCalling TestMakeEvent(\n  %v\n),\tresult: %d\n", testCase.actual, resp.GetEventId())

		// Assert:
		if err != nil {
			t.Errorf("TestMakeEvent(%v) got unexpected error", testCase.actual)
		}
		assert.Equal(t, testCase.expected.ID, resp.GetEventId())
	}
}

func TestGetEvent(t *testing.T) {
	// Arrange:
	s := service.RunEventsService()

	testTable := []struct {
		actual   models.Events
		expected models.Events
	}{
		{
			actual: models.Events{
				SenderID: 0,
				ID:       0,
			},
			expected: models.Events{
				SenderID: 0,
				ID:       0,
				Time:     0,
				Name:     "",
			},
		},
	}

	// Act:
	for _, testCase := range testTable {
		req := &eventmanager.GetEventRequest{
			EventId: testCase.actual.ID,
		}
		resp, err := s.GetEvent(context.Background(), req)
		t.Logf("\nCalling TestGetEvent(\n  %v\n),\tresult: %d\n", testCase.actual, resp.GetEventId())

		// Assert:
		if err != nil {
			assert.Equal(t, err, service.ErrEventNotFound)
		} else {
			assert.Equal(t, testCase.expected.SenderID, resp.GetSenderId())
			assert.Equal(t, testCase.expected.ID, resp.GetEventId())
			assert.Equal(t, testCase.expected.Time, resp.GetTime())
			assert.Equal(t, testCase.expected.Name, resp.GetName())
		}
	}
}
