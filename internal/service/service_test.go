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
	s := service.Server{
		EventsByClient: make(map[int64]map[int64]models.Events),
	}
	currentTime := time.Now().Local().UnixMilli()

	testTable := []struct {
		actual   models.Events
		expected models.Events
	}{
		{
			actual: models.Events{
				SenderID: 1,
				Time:     currentTime,
				Name:     "TestUser1",
			},
			expected: models.Events{
				ID: 1,
			},
		},
		{
			actual: models.Events{
				SenderID: 1,
				Time:     currentTime,
				Name:     "TestUser1Again",
			},
			expected: models.Events{
				ID: 2,
			},
		},
		{
			actual: models.Events{
				SenderID: 2,
				Time:     currentTime,
				Name:     "TestUser2",
			},
			expected: models.Events{
				ID: 1,
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
		if err != nil {
			t.Errorf("TestMakeEvent(%v) got unexpected error", testCase.actual)
		}
		t.Logf("Calling TestMakeEvent(%v), result: %d\n", testCase.actual, resp.GetEventId())
		// Assert:
		assert.Equal(t, testCase.expected.ID, resp.GetEventId())
	}
}
