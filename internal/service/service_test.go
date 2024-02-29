package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/Raimguzhinov/simple-grpc/internal/models"
	"github.com/Raimguzhinov/simple-grpc/internal/service"
	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

func TestRunEventsService(t *testing.T) {
	// Arrange:
	uninitializedServer := service.Server{}
	testCases := []struct {
		name     string
		actual   bool
		expected bool
	}{
		{
			name:     "Test 1",
			actual:   service.RunEventsService().IsInitialized(),
			expected: true,
		},
		{
			name:     "Test 2 (Negative Case)",
			actual:   uninitializedServer.IsInitialized(),
			expected: false,
		},
	}

	// Act:
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Assert:
			assert.Equal(t, testCase.expected, testCase.actual)
		})
	}
}

func TestMakeEvent(t *testing.T) {
	// Arrange:
	s := service.RunEventsService()
	currentTime := time.Now().UnixMilli()

	testTable := []struct {
		name     string
		actual   models.Events
		expected models.Events
	}{
		{
			name: "Test 1",
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
			name: "Test 2",
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
			name: "Test 3",
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
			name: "Test 4",
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
			name: "Test 5",
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
		t.Run(testCase.name, func(t *testing.T) {
			req := &eventmanager.MakeEventRequest{
				SenderId: testCase.actual.SenderID,
				Time:     testCase.actual.Time,
				Name:     testCase.actual.Name,
			}
			resp, err := s.MakeEvent(context.Background(), req)
			t.Logf("\nCalling TestMakeEvent(\n  %v\n),\nresult: %v\n\n", testCase.actual, resp.GetEventId())

			// Assert:
			if err != nil {
				t.Errorf("TestMakeEvent(%v) got unexpected error", testCase.actual)
			}
			assert.Equal(t, testCase.expected.ID, resp.GetEventId())
		})
	}
}

func mockEventMaker(s *service.Server, currentTime int64) error {
	oneMonthLater := time.UnixMilli(currentTime).AddDate(0, 1, 0).UnixMilli()
	eventsData := []struct {
		SenderId int64
		Time     int64
		Name     string
	}{
		{SenderId: 1, Time: currentTime, Name: "User1"},
		{SenderId: 2, Time: currentTime, Name: "User2"},
		{SenderId: 2, Time: oneMonthLater, Name: "User2Again"},
		{SenderId: 3, Time: currentTime, Name: "User3"},
	}

	for _, eventData := range eventsData {
		_, err := s.MakeEvent(context.Background(), &eventmanager.MakeEventRequest{
			SenderId: eventData.SenderId,
			Time:     eventData.Time,
			Name:     eventData.Name,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func setupTest(t *testing.T, currentTime time.Time) *service.Server {
	s := service.RunEventsService()
	t.Run("Mocking events", func(t *testing.T) {
		err := mockEventMaker(s, currentTime.UnixMilli())
		if err != nil {
			t.Error(err)
		}
	})
	t.Log("Events created")
	return s
}

func TestGetEvent(t *testing.T) {
	// Arrange:
	currentTimeT := time.Now()
	s := setupTest(t, currentTimeT)
	currentTime := currentTimeT.UnixMilli()
	loc, _ := time.LoadLocation("America/New_York")

	testTable := []struct {
		name     string
		actual   models.Events
		expected models.Events
	}{
		{
			name: "Test 1",
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
		{
			name: "Test 2",
			actual: models.Events{
				SenderID: 1,
				ID:       1,
			},
			expected: models.Events{
				SenderID: 1,
				ID:       1,
				Time:     currentTime,
				Name:     "User1",
			},
		},
		{
			name: "Test 3",
			actual: models.Events{
				SenderID: 2,
				ID:       1,
			},
			expected: models.Events{
				SenderID: 2,
				ID:       1,
				Time:     currentTime,
				Name:     "User2",
			},
		},
		{
			name: "Test 4",
			actual: models.Events{
				SenderID: 1,
				ID:       2,
			},
			expected: models.Events{
				SenderID: 0,
				ID:       0,
				Time:     0,
				Name:     "",
			},
		},
		{
			name: "Test 5",
			actual: models.Events{
				SenderID: 999,
				ID:       1,
			},
			expected: models.Events{
				SenderID: 0,
				ID:       0,
				Time:     0,
				Name:     "",
			},
		},
		{
			name: "Test 6",
			actual: models.Events{
				SenderID: 3,
				ID:       1,
			},
			expected: models.Events{
				SenderID: 3,
				ID:       1,
				Time:     time.UnixMilli(currentTime).In(loc).UnixMilli(),
				Name:     "User3",
			},
		},
	}

	// Act:
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req := &eventmanager.GetEventRequest{
				SenderId: testCase.actual.SenderID,
				EventId:  testCase.actual.ID,
			}
			resp, err := s.GetEvent(context.Background(), req)
			t.Logf("\nCalling TestGetEvent(\n  %v\n),\nresult: %v\n\n", testCase.actual, resp.GetEventId())

			// Assert:
			if err != nil {
				assert.Equal(t, err, service.ErrEventNotFound)
			}
			assert.Equal(t, testCase.expected.SenderID, resp.GetSenderId())
			assert.Equal(t, testCase.expected.ID, resp.GetEventId())
			assert.Equal(t, testCase.expected.Time, resp.GetTime())
			assert.Equal(t, testCase.expected.Name, resp.GetName())
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	// Arrange:
	currentTimeT := time.Now()
	s := setupTest(t, currentTimeT)

	testTable := []struct {
		name     string
		actual   models.Events
		expected models.Events
	}{
		{
			name: "Test 1",
			actual: models.Events{
				SenderID: 0,
				ID:       0,
			},
			expected: models.Events{
				ID: 0,
			},
		},
		{
			name: "Test 2",
			actual: models.Events{
				SenderID: 1,
				ID:       1,
			},
			expected: models.Events{
				ID: 1,
			},
		},
		{
			name: "Test 3",
			actual: models.Events{
				SenderID: 2,
				ID:       1,
			},
			expected: models.Events{
				ID: 1,
			},
		},
		{
			name: "Test 4",
			actual: models.Events{
				SenderID: 1,
				ID:       2,
			},
			expected: models.Events{
				ID: 0,
			},
		},
		{
			name: "Test 5",
			actual: models.Events{
				SenderID: 999,
				ID:       1,
			},
			expected: models.Events{
				ID: 0,
			},
		},
	}

	// Act:
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req := &eventmanager.DeleteEventRequest{
				SenderId: testCase.actual.SenderID,
				EventId:  testCase.actual.ID,
			}
			resp, err := s.DeleteEvent(context.Background(), req)
			t.Logf("\nCalling TestDeleteEvent(\n  %v\n),\nresult: %v\n\n", testCase.actual, resp.GetEventId())

			// Assert:
			if err != nil {
				assert.Equal(t, err, service.ErrEventNotFound)
			}
			assert.Equal(t, testCase.expected.ID, resp.GetEventId())
		})
	}
}

func TestGetEvents(t *testing.T) {
	// Arrange:
	currentTime := time.Now()
	s := setupTest(t, currentTime)
	oneYearAgo := currentTime.AddDate(-1, 0, 0).UnixMilli()
	oneYearLater := currentTime.AddDate(1, 0, 0).UnixMilli()

	testTable := []struct {
		name      string
		actual    models.Events
		timerange [2]int64
		expected  []models.Events
	}{
		{
			name: "Test 1",
			actual: models.Events{
				SenderID: 0,
			},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected:  []models.Events{},
		},
		{
			name: "Test 2",
			actual: models.Events{
				SenderID: 1,
			},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected: []models.Events{
				{
					SenderID: 1,
					ID:       1,
					Time:     currentTime.UnixMilli(),
					Name:     "User1",
				},
			},
		},
		{
			name: "Test 3",
			actual: models.Events{
				SenderID: 2,
			},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected: []models.Events{
				{
					SenderID: 2,
					ID:       1,
					Time:     currentTime.UnixMilli(),
					Name:     "User2",
				},
				{
					SenderID: 2,
					ID:       2,
					Time:     currentTime.AddDate(0, 1, 0).UnixMilli(),
					Name:     "User2Again",
				},
			},
		},
		{
			name: "Test 4",
			actual: models.Events{
				SenderID: 3,
			},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected: []models.Events{
				{
					SenderID: 3,
					ID:       1,
					Time:     currentTime.UnixMilli(),
					Name:     "User3",
				},
			},
		},
		{
			name: "Test 5",
			actual: models.Events{
				SenderID: 3,
			},
			timerange: [2]int64{oneYearAgo, oneYearAgo},
			expected:  []models.Events{},
		},
		{
			name: "Test 6",
			actual: models.Events{
				SenderID: 3,
			},
			timerange: [2]int64{oneYearLater, oneYearLater},
			expected:  []models.Events{},
		},
	}
	// Act:
	for _, testCase := range testTable {
		mockStream := &mockEventStream{}
		t.Run(testCase.name, func(t *testing.T) {
			req := &eventmanager.GetEventsRequest{
				SenderId: testCase.actual.SenderID,
				FromTime: testCase.timerange[0],
				ToTime:   testCase.timerange[1],
			}
			mockStream = &mockEventStream{
				sentEvents: make([]*eventmanager.Event, 0),
			}
			err := s.GetEvents(req, mockStream)
			if err != nil {
				t.Logf("\nCalling TestGetEvents(\n  %v\n),\nresult: %v\n\n", testCase.actual, err)
				assert.Equal(t, err, service.ErrEventNotFound)
				t.Skipf("Skipping TestGetEvents: %v", err)
			}
			t.Log(mockStream.sentEvents)
			for i, resp := range mockStream.sentEvents {
				t.Logf("\nCalling TestGetEvents(\n  %v\n),\nresult: %v\n\n", testCase.actual, resp)
				// Assert:
				assert.Equal(t, service.AccumulateEvent(testCase.expected[i]).GetSenderId(), resp.GetSenderId())
				assert.Equal(t, service.AccumulateEvent(testCase.expected[i]).GetEventId(), resp.GetEventId())
				assert.Equal(t, service.AccumulateEvent(testCase.expected[i]).GetName(), resp.GetName())
				assert.Equal(t, service.AccumulateEvent(testCase.expected[i]).GetTime(), resp.GetTime())
			}
		})
	}
}

type mockEventStream struct {
	sentEvents []*eventmanager.Event
	grpc.ServerStream
}

func (m *mockEventStream) Send(event *eventmanager.Event) error {
	m.sentEvents = append(m.sentEvents, event)
	return nil
}

func TestAccumulateEvent(t *testing.T) {
	// Arrange:
	testCases := []struct {
		name     string
		event    models.Events
		expected *eventmanager.Event
	}{
		{
			name: "Test1",
			event: models.Events{
				SenderID: 1,
				ID:       1,
				Time:     123456789,
				Name:     "TestEvent",
			},
			expected: &eventmanager.Event{
				SenderId: 1,
				EventId:  1,
				Time:     123456789,
				Name:     "TestEvent",
			},
		},
		{
			name:     "Test2",
			event:    models.Events{},
			expected: &eventmanager.Event{},
		},
	}

	// Act:
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := service.AccumulateEvent(testCase.event)
			// Assert:
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
