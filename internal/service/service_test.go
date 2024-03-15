package service_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/Raimguzhinov/simple-grpc/internal/models"
	"github.com/Raimguzhinov/simple-grpc/internal/service"
	eventcrtl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type testStruct struct {
	name     string
	actual   models.Event
	expected models.Event
}

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
	t.Parallel()
	// Arrange:
	s := service.RunEventsService()
	oneMonthLater := time.Now().AddDate(0, 1, 0).UnixMilli()

	var testTable []testStruct
	for i := 0; i < 100; i++ {
		testTable = append(testTable, testStruct{
			name: fmt.Sprintf("Test %d", i),
			actual: models.Event{
				SenderID: int64(i),
				Time:     oneMonthLater,
				Name:     fmt.Sprintf("User %d", i),
			},
			expected: models.Event{
				ID: 1,
			},
		})
	}

	var wg sync.WaitGroup
	type testResult struct {
		resp     *eventcrtl.EventIdAvail
		err      error
		expected models.Event
	}
	resultCh := make(chan testResult, 1)

	// Act:
	for _, testCase := range testTable {
		wg.Add(1)
		go func(testCase testStruct) {
			defer wg.Done()

			req := &eventcrtl.MakeEventRequest{
				SenderId: testCase.actual.SenderID,
				Time:     testCase.actual.Time,
				Name:     testCase.actual.Name,
			}
			resp, err := s.MakeEvent(context.Background(), req)
			if err != nil {
				err = fmt.Errorf("TestMakeEvent(%v) got unexpected error", testCase.actual)
			}
			t.Logf(
				"\nCalling TestMakeEvent(\n  %v\n),\nresult: %v\n\n",
				testCase.actual,
				resp.GetEventId(),
			)
			resultCh <- testResult{resp, err, testCase.expected}
		}(testCase)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Assert:
	for res := range resultCh {
		assert.Nil(t, res.err)
		assert.Equal(t, res.expected.ID, res.resp.GetEventId())
	}
}

func MockEventMaker(t *testing.T, s *service.Server, givenTime int64) error {
	oneMonthLater := time.UnixMilli(givenTime).AddDate(0, 1, 0).UnixMilli()
	eventsData := []struct {
		SenderId int64
		Time     int64
		Name     string
	}{
		{SenderId: 1, Time: givenTime, Name: "User1"},
		{SenderId: 2, Time: givenTime, Name: "User2"},
		{SenderId: 2, Time: oneMonthLater, Name: "User2Again"},
		{SenderId: 3, Time: givenTime, Name: "User3"},
	}

	for _, eventData := range eventsData {
		resp, err := s.MakeEvent(context.Background(), &eventcrtl.MakeEventRequest{
			SenderId: eventData.SenderId,
			Time:     eventData.Time,
			Name:     eventData.Name,
		})
		t.Log(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetupTest(t *testing.T, currentTime time.Time) *service.Server {
	s := service.RunEventsService()
	t.Run("Mocking events", func(t *testing.T) {
		err := MockEventMaker(t, s, currentTime.UnixMilli())
		if err != nil {
			t.Error(err)
		}
	})
	t.Log("Events created")
	return s
}

func TestGetEvent(t *testing.T) {
	t.Parallel()
	// Arrange:
	oneMonthLaterT := time.Now().AddDate(0, 1, 0)
	s := SetupTest(t, oneMonthLaterT)
	oneMonthLater := oneMonthLaterT.UnixMilli()
	loc, _ := time.LoadLocation("America/New_York")

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name: "Test 1",
			actual: models.Event{
				SenderID: 0,
				ID:       0,
			},
			expected: models.Event{
				SenderID: 0,
				ID:       0,
				Time:     0,
				Name:     "",
			},
		},
		testStruct{
			name: "Test 2",
			actual: models.Event{
				SenderID: 1,
				ID:       1,
			},
			expected: models.Event{
				SenderID: 1,
				ID:       1,
				Time:     oneMonthLater,
				Name:     "User1",
			},
		},
		testStruct{
			name: "Test 3",
			actual: models.Event{
				SenderID: 2,
				ID:       1,
			},
			expected: models.Event{
				SenderID: 2,
				ID:       1,
				Time:     oneMonthLater,
				Name:     "User2",
			},
		},
		testStruct{
			name: "Test 4",
			actual: models.Event{
				SenderID: 1,
				ID:       2,
			},
			expected: models.Event{
				SenderID: 0,
				ID:       0,
				Time:     0,
				Name:     "",
			},
		},
		testStruct{
			name: "Test 5",
			actual: models.Event{
				SenderID: 999,
				ID:       1,
			},
			expected: models.Event{
				SenderID: 0,
				ID:       0,
				Time:     0,
				Name:     "",
			},
		},
		testStruct{
			name: "Test 6",
			actual: models.Event{
				SenderID: 3,
				ID:       1,
			},
			expected: models.Event{
				SenderID: 3,
				ID:       1,
				Time:     time.UnixMilli(oneMonthLater).In(loc).UnixMilli(),
				Name:     "User3",
			},
		},
	)

	var wg sync.WaitGroup
	type testResult struct {
		resp     *eventcrtl.Event
		err      error
		expected models.Event
	}
	resultCh := make(chan testResult, 1)

	// Act:
	for _, testCase := range testTable {
		wg.Add(1)
		go func(testCase testStruct) {
			defer wg.Done()

			req := &eventcrtl.GetEventRequest{
				SenderId: testCase.actual.SenderID,
				EventId:  testCase.actual.ID,
			}
			resp, err := s.GetEvent(context.Background(), req)
			t.Logf(
				"\nCalling TestGetEvent(\n  %v\n),\nresult: %v\n\n",
				testCase.actual,
				resp.GetEventId(),
			)

			resultCh <- testResult{resp, err, testCase.expected}
		}(testCase)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Assert:
	for res := range resultCh {
		if res.err != nil {
			assert.Equal(t, service.ErrEventNotFound, res.err)
		}
		assert.Equal(t, res.expected.SenderID, res.resp.GetSenderId())
		assert.Equal(t, res.expected.ID, res.resp.GetEventId())
		assert.Equal(t, res.expected.Time, res.resp.GetTime())
		assert.Equal(t, res.expected.Name, res.resp.GetName())
	}
}

func TestDeleteEvent(t *testing.T) {
	t.Parallel()
	// Arrange:
	oneMonthLaterT := time.Now().AddDate(0, 1, 0)
	s := SetupTest(t, oneMonthLaterT)

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name: "Test 1",
			actual: models.Event{
				SenderID: 0,
				ID:       0,
			},
			expected: models.Event{
				ID: 0,
			},
		},
		testStruct{
			name: "Test 2",
			actual: models.Event{
				SenderID: 1,
				ID:       1,
			},
			expected: models.Event{
				ID: 1,
			},
		},
		testStruct{
			name: "Test 3",
			actual: models.Event{
				SenderID: 2,
				ID:       1,
			},
			expected: models.Event{
				ID: 1,
			},
		},
		testStruct{
			name: "Test 4",
			actual: models.Event{
				SenderID: 1,
				ID:       2,
			},
			expected: models.Event{
				ID: 0,
			},
		},
		testStruct{
			name: "Test 5",
			actual: models.Event{
				SenderID: 999,
				ID:       1,
			},
			expected: models.Event{
				ID: 0,
			},
		},
	)

	var wg sync.WaitGroup
	type testResult struct {
		resp     *eventcrtl.EventIdAvail
		err      error
		expected models.Event
	}
	resultCh := make(chan testResult, 1)

	// Act:
	for _, testCase := range testTable {
		wg.Add(1)
		go func(testCase testStruct) {
			defer wg.Done()

			req := &eventcrtl.DeleteEventRequest{
				SenderId: testCase.actual.SenderID,
				EventId:  testCase.actual.ID,
			}
			resp, err := s.DeleteEvent(context.Background(), req)
			t.Logf(
				"\nCalling TestDeleteEvent(\n  %v\n),\nresult: %v\n\n",
				testCase.actual,
				resp.GetEventId(),
			)
			resultCh <- testResult{resp, err, testCase.expected}
		}(testCase)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Assert:
	for res := range resultCh {
		if res.err != nil {
			assert.Equal(t, service.ErrEventNotFound, res.err)
		}
		assert.Equal(t, res.expected.ID, res.resp.GetEventId())
	}
}

func TestGetEvents(t *testing.T) {
	t.Parallel()
	// Arrange:
	oneMonthLaterT := time.Now().AddDate(0, 1, 0)
	s := SetupTest(t, oneMonthLaterT)
	oneMonthLater := oneMonthLaterT.UnixMilli()
	oneYearAgo := time.Now().AddDate(-1, 0, 0).UnixMilli()
	oneYearLater := time.Now().AddDate(1, 0, 0).UnixMilli()

	type testStruct struct {
		name      string
		actual    models.Event
		timerange [2]int64
		expected  []models.Event
	}

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name: "Test 1",
			actual: models.Event{
				SenderID: 0,
			},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected:  []models.Event{},
		},
		testStruct{
			name: "Test 2",
			actual: models.Event{
				SenderID: 1,
			},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected: []models.Event{
				{
					SenderID: 1,
					ID:       1,
					Time:     oneMonthLater,
					Name:     "User1",
				},
			},
		},
		testStruct{
			name: "Test 3",
			actual: models.Event{
				SenderID: 2,
			},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected: []models.Event{
				{
					SenderID: 2,
					ID:       1,
					Time:     oneMonthLater,
					Name:     "User2",
				},
				{
					SenderID: 2,
					ID:       2,
					Time:     oneMonthLaterT.AddDate(0, 1, 0).UnixMilli(),
					Name:     "User2Again",
				},
			},
		},
		testStruct{
			name: "Test 4",
			actual: models.Event{
				SenderID: 3,
			},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected: []models.Event{
				{
					SenderID: 3,
					ID:       1,
					Time:     oneMonthLater,
					Name:     "User3",
				},
			},
		},
		testStruct{
			name: "Test 5",
			actual: models.Event{
				SenderID: 3,
			},
			timerange: [2]int64{oneYearAgo, oneYearAgo},
			expected:  []models.Event{},
		},
		testStruct{
			name: "Test 6",
			actual: models.Event{
				SenderID: 3,
			},
			timerange: [2]int64{oneYearLater, oneYearLater},
			expected:  []models.Event{},
		},
	)

	var wg sync.WaitGroup
	type testResult struct {
		resp     *eventcrtl.Event
		err      error
		expected *eventcrtl.Event
	}
	resultCh := make(chan testResult, 1)
	_ = resultCh

	// Act:
	for _, testCase := range testTable {
		wg.Add(1)
		mockStream := &mockEventStream{}
		go func(testCase testStruct) {
			defer wg.Done()

			req := &eventcrtl.GetEventsRequest{
				SenderId: testCase.actual.SenderID,
				FromTime: testCase.timerange[0],
				ToTime:   testCase.timerange[1],
			}
			mockStream = &mockEventStream{
				SentEvents: make([]*eventcrtl.Event, 0),
			}
			err := s.GetEvents(req, mockStream)
			if err != nil {
				t.Logf("\nCalling TestGetEvents(\n  %v\n),\nresult: %v\n\n", testCase.actual, err)
				assert.Equal(t, err, service.ErrEventNotFound)
				return
			}
			t.Log(mockStream.SentEvents)

			for i, resp := range mockStream.SentEvents {
				t.Logf("\nCalling TestGetEvents(\n  %v\n),\nresult: %v\n\n", testCase.actual, resp)
				resultCh <- testResult{resp, err, service.AccumulateEvent(testCase.expected[i])}
			}
		}(testCase)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Assert:
	for res := range resultCh {
		assert.Equal(t, res.expected.GetSenderId(), res.resp.GetSenderId())
		assert.Equal(t, res.expected.GetEventId(), res.resp.GetEventId())
		assert.Equal(t, res.expected.GetName(), res.resp.GetName())
		assert.Equal(t, res.expected.GetTime(), res.resp.GetTime())
	}
}

type mockEventStream struct {
	SentEvents []*eventcrtl.Event
	grpc.ServerStream
}

func (m *mockEventStream) Send(event *eventcrtl.Event) error {
	m.SentEvents = append(m.SentEvents, event)
	return nil
}

func TestAccumulateEvent(t *testing.T) {
	// Arrange:
	testCases := []struct {
		name     string
		event    models.Event
		expected *eventcrtl.Event
	}{
		{
			name: "Test1",
			event: models.Event{
				SenderID: 1,
				ID:       1,
				Time:     123456789,
				Name:     "TestEvent",
			},
			expected: &eventcrtl.Event{
				SenderId: 1,
				EventId:  1,
				Time:     123456789,
				Name:     "TestEvent",
			},
		},
		{
			name:     "Test2",
			event:    models.Event{},
			expected: &eventcrtl.Event{},
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
