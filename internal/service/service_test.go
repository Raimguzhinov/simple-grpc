package service_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/Raimguzhinov/simple-grpc/internal/models"
	"github.com/Raimguzhinov/simple-grpc/internal/service"
	eventcrtl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
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
			require.Equal(t, testCase.expected, testCase.actual)
		})
	}
}

func TestMakeEvent(t *testing.T) {
	t.Parallel()
	// Arrange:
	s := service.RunEventsService()
	oneMonthLater := time.Now().AddDate(0, 1, 0).UnixMilli()

	type testStruct struct {
		name     string
		actual   models.Event
		expected bool
	}
	var testTable []testStruct
	for i := 0; i < 100; i++ {
		testTable = append(testTable, testStruct{
			name: fmt.Sprintf("Test %d", i),
			actual: models.Event{
				SenderID: int64(i),
				Time:     oneMonthLater,
				Name:     fmt.Sprintf("User %d", i),
			},
			expected: true,
		})
	}

	var wg sync.WaitGroup
	type testResult struct {
		actual   bool
		err      error
		expected bool
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
			actual := s.IsEventExist(testCase.actual.SenderID, resp.GetEventId())
			resultCh <- testResult{actual, err, testCase.expected}
		}(testCase)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Assert:
	for res := range resultCh {
		require.Nil(t, res.err)
		require.Equal(t, res.expected, res.actual)
	}
}

func TestMakeEventMultiSession(t *testing.T) {
	t.Parallel()
	// Arrange:
	s := service.RunEventsService()
	oneMonthLater := time.Now().AddDate(0, 1, 0).UnixMilli()

	type testStruct struct {
		name   string
		actual models.Event
	}
	var testTable []testStruct
	expectedEvents := 2

	for i := 0; i < 100; i += expectedEvents {
		testTable = append(testTable,
			testStruct{
				name: fmt.Sprintf("Test %d", i),
				actual: models.Event{
					SenderID: int64(i),
					Time:     oneMonthLater,
					Name:     fmt.Sprintf("User %d Maybe First", i),
				},
			},
			testStruct{
				name: fmt.Sprintf("Test %d", i+1),
				actual: models.Event{
					SenderID: int64(i),
					Time:     oneMonthLater,
					Name:     fmt.Sprintf("User %d Maybe Again", i),
				},
			},
		)
	}

	var wg sync.WaitGroup
	type testResult struct {
		actual   bool
		err      error
		expected bool
	}
	resultCh := make(chan testResult, 1)
	recurringEvents := make(map[int64]bool)
	var mu sync.Mutex

	// Act:
	for _, testCase := range testTable {
		wg.Add(1)
		go func(testCase testStruct) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()

			req := &eventcrtl.MakeEventRequest{
				SenderId: testCase.actual.SenderID,
				Time:     testCase.actual.Time,
				Name:     testCase.actual.Name,
			}
			resp, err := s.MakeEvent(context.Background(), req)
			if err != nil {
				err = fmt.Errorf(
					"TestMakeEventMultiSession(%v) got unexpected error",
					testCase.actual,
				)
			}
			t.Logf(
				"\nCalling TestMakeEventMultiSession(\n  %v\n),\nresult: %v\n\n",
				testCase.actual,
				resp.GetEventId(),
			)

			senderID := testCase.actual.SenderID
			if _, ok := recurringEvents[senderID]; ok {
				actual := s.IsEventExist(senderID, resp.GetEventId())
				resultCh <- testResult{actual, err, true}
			} else {
				recurringEvents[senderID] = true
			}
		}(testCase)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Assert:
	for res := range resultCh {
		require.Nil(t, res.err)
		require.Equal(t, res.expected, res.actual)
	}
}

func MockEventMaker(
	t *testing.T,
	s *service.Server,
	givenTime int64,
) (*map[int64]map[string][]byte, error) {
	eventsIdMap := make(map[int64]map[string][]byte)
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
			return nil, err
		}
		if _, ok := eventsIdMap[eventData.SenderId]; !ok {
			eventsIdMap[eventData.SenderId] = make(map[string][]byte)
		}
		eventsIdMap[eventData.SenderId][eventData.Name] = resp.EventId
	}
	return &eventsIdMap, nil
}

func SetupTest(
	t *testing.T,
	currentTime time.Time,
) (*service.Server, *map[int64]map[string][]byte) {
	s := service.RunEventsService()
	var idMap *map[int64]map[string][]byte
	t.Run("Mocking events", func(t *testing.T) {
		var err error
		idMap, err = MockEventMaker(t, s, currentTime.UnixMilli())
		if err != nil {
			t.Error(err)
		}
	})
	t.Log("Events created")
	return s, idMap
}

func TestGetEvent(t *testing.T) {
	t.Parallel()
	// Arrange:
	oneMonthLaterT := time.Now().AddDate(0, 1, 0)
	s, idMapPtr := SetupTest(t, oneMonthLaterT)
	idMap := *idMapPtr
	oneMonthLater := oneMonthLaterT.UnixMilli()
	loc, _ := time.LoadLocation("America/New_York")

	type testStruct struct {
		name        string
		actual      *eventcrtl.GetEventRequest
		expected    *eventcrtl.Event
		expectedErr error
	}
	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name: "Test 1",
			actual: &eventcrtl.GetEventRequest{
				SenderId: 0,
				EventId:  []byte("0000-0000-0000"),
			},
			expected:    nil,
			expectedErr: service.ErrEventNotFound,
		},
		testStruct{
			name: "Test 2",
			actual: &eventcrtl.GetEventRequest{
				SenderId: 1,
				EventId:  idMap[1]["User1"],
			},
			expected: &eventcrtl.Event{
				SenderId: 1,
				EventId:  idMap[1]["User1"],
				Time:     oneMonthLater,
				Name:     "User1",
			},
			expectedErr: nil,
		},
		testStruct{
			name: "Test 3",
			actual: &eventcrtl.GetEventRequest{
				SenderId: 2,
				EventId:  idMap[2]["User2"],
			},
			expected: &eventcrtl.Event{
				SenderId: 2,
				EventId:  idMap[2]["User2"],
				Time:     oneMonthLater,
				Name:     "User2",
			},
			expectedErr: nil,
		},
		testStruct{
			name: "Test 4",
			actual: &eventcrtl.GetEventRequest{
				SenderId: 3,
				EventId:  idMap[1]["User1"],
			},
			expected:    nil,
			expectedErr: service.ErrEventNotFound,
		},
		testStruct{
			name: "Test 4",
			actual: &eventcrtl.GetEventRequest{
				SenderId: 999,
				EventId:  idMap[2]["User2"],
			},
			expected:    nil,
			expectedErr: service.ErrEventNotFound,
		},
		testStruct{
			name: "Test 6",
			actual: &eventcrtl.GetEventRequest{
				SenderId: 3,
				EventId:  idMap[3]["User3"],
			},
			expected: &eventcrtl.Event{
				SenderId: 3,
				EventId:  idMap[3]["User3"],
				Time:     time.UnixMilli(oneMonthLater).In(loc).UnixMilli(),
				Name:     "User3",
			},
			expectedErr: nil,
		},
	)

	var wg sync.WaitGroup
	type testResult struct {
		resp, expected   *eventcrtl.Event
		err, expectedErr error
	}
	resultCh := make(chan testResult, 1)

	// Act:
	for _, testCase := range testTable {
		wg.Add(1)
		go func(testCase testStruct) {
			defer wg.Done()

			req := testCase.actual
			resp, err := s.GetEvent(context.Background(), req)
			t.Logf(
				"\nCalling TestGetEvent(\n  %v\n),\nresult: %v\n\n",
				testCase.actual,
				resp.GetEventId(),
			)

			resultCh <- testResult{resp, testCase.expected, err, testCase.expectedErr}
		}(testCase)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Assert:
	for res := range resultCh {
		if res.expectedErr == nil {
			require.NoError(t, res.err)
		} else {
			require.EqualError(t, res.err, res.expectedErr.Error())
		}

		if res.expected == nil {
			require.Nil(t, res.resp)
		} else {
			require.Equal(t, res.expected.SenderId, res.resp.GetSenderId())
			require.Equal(t, res.expected.EventId, res.resp.GetEventId())
			require.Equal(t, res.expected.Time, res.resp.GetTime())
			require.Equal(t, res.expected.Name, res.resp.GetName())
		}
	}
}

func TestDeleteEvent(t *testing.T) {
	t.Parallel()
	// Arrange:
	oneMonthLaterT := time.Now().AddDate(0, 1, 0)
	s, idMapPtr := SetupTest(t, oneMonthLaterT)
	idMap := *idMapPtr

	type testStruct struct {
		name        string
		actual      *eventcrtl.DeleteEventRequest
		expected    *eventcrtl.EventIdAvail
		expectedErr error
	}
	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name: "Test 1",
			actual: &eventcrtl.DeleteEventRequest{
				SenderId: 0,
				EventId:  []byte("0000-0000-0000"),
			},
			expected:    nil,
			expectedErr: service.ErrEventNotFound,
		},
		testStruct{
			name: "Test 2",
			actual: &eventcrtl.DeleteEventRequest{
				SenderId: 1,
				EventId:  idMap[1]["User1"],
			},
			expected:    &eventcrtl.EventIdAvail{EventId: idMap[1]["User1"]},
			expectedErr: nil,
		},
		testStruct{
			name: "Test 3",
			actual: &eventcrtl.DeleteEventRequest{
				SenderId: 2,
				EventId:  idMap[2]["User2"],
			},
			expected:    &eventcrtl.EventIdAvail{EventId: idMap[2]["User2"]},
			expectedErr: nil,
		},
		testStruct{
			name: "Test 4",
			actual: &eventcrtl.DeleteEventRequest{
				SenderId: 3,
				EventId:  []byte("999-9999-9999"),
			},
			expected:    nil,
			expectedErr: service.ErrEventNotFound,
		},
		testStruct{
			name: "Test 5",
			actual: &eventcrtl.DeleteEventRequest{
				SenderId: 999,
				EventId:  idMap[1]["User1"],
			},
			expected:    nil,
			expectedErr: service.ErrEventNotFound,
		},
	)

	var wg sync.WaitGroup
	type testResult struct {
		resp, expected   *eventcrtl.EventIdAvail
		err, expectedErr error
	}
	resultCh := make(chan testResult, 1)

	// Act:
	for _, testCase := range testTable {
		wg.Add(1)
		go func(testCase testStruct) {
			defer wg.Done()

			req := testCase.actual
			resp, err := s.DeleteEvent(context.Background(), req)
			t.Logf(
				"\nCalling TestDeleteEvent(\n  %v\n),\nresult: %v\n\n",
				testCase.actual,
				resp.GetEventId(),
			)
			resultCh <- testResult{resp, testCase.expected, err, testCase.expectedErr}
		}(testCase)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Assert:
	for res := range resultCh {
		if res.expectedErr == nil {
			require.NoError(t, res.err)
		} else {
			require.EqualError(t, res.err, res.expectedErr.Error())
		}

		if res.expected == nil {
			require.Nil(t, res.resp)
		} else {
			require.Equal(t, res.expected.EventId, res.resp.GetEventId())
		}
	}
}

func TestGetEvents(t *testing.T) {
	t.Parallel()
	// Arrange:
	oneMonthLaterT := time.Now().AddDate(0, 1, 0)
	s, idMapPtr := SetupTest(t, oneMonthLaterT)
	idMap := *idMapPtr
	oneMonthLater := oneMonthLaterT.UnixMilli()
	oneYearAgo := time.Now().AddDate(-1, 0, 0).UnixMilli()
	oneYearLater := time.Now().AddDate(1, 0, 0).UnixMilli()

	type testStruct struct {
		name        string
		actual      models.Event
		timerange   [2]int64
		expected    []*eventcrtl.Event
		expectedErr error
	}

	var testTable []testStruct
	testTable = append(testTable,
		testStruct{
			name:        "Test 1",
			actual:      models.Event{SenderID: 0},
			timerange:   [2]int64{oneYearAgo, oneYearLater},
			expected:    []*eventcrtl.Event{},
			expectedErr: service.ErrEventNotFound,
		},
		testStruct{
			name:      "Test 2",
			actual:    models.Event{SenderID: 1},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected: []*eventcrtl.Event{
				{
					SenderId: 1,
					EventId:  idMap[1]["User1"],
					Time:     oneMonthLater,
					Name:     "User1",
				},
			},
			expectedErr: nil,
		},
		testStruct{
			name:      "Test 3",
			actual:    models.Event{SenderID: 2},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected: []*eventcrtl.Event{
				{
					SenderId: 2,
					EventId:  idMap[2]["User2"],
					Time:     oneMonthLater,
					Name:     "User2",
				},
				{
					SenderId: 2,
					EventId:  idMap[2]["User2Again"],
					Time:     oneMonthLaterT.AddDate(0, 1, 0).UnixMilli(),
					Name:     "User2Again",
				},
			},
			expectedErr: nil,
		},
		testStruct{
			name:      "Test 4",
			actual:    models.Event{SenderID: 3},
			timerange: [2]int64{oneYearAgo, oneYearLater},
			expected: []*eventcrtl.Event{
				{
					SenderId: 3,
					EventId:  idMap[3]["User3"],
					Time:     oneMonthLater,
					Name:     "User3",
				},
			},
			expectedErr: nil,
		},
		testStruct{
			name:        "Test 5",
			actual:      models.Event{SenderID: 3},
			timerange:   [2]int64{oneYearAgo, oneYearAgo},
			expected:    []*eventcrtl.Event{},
			expectedErr: service.ErrEventNotFound,
		},
		testStruct{
			name:        "Test 6",
			actual:      models.Event{SenderID: 3},
			timerange:   [2]int64{oneYearLater, oneYearLater},
			expected:    []*eventcrtl.Event{},
			expectedErr: service.ErrEventNotFound,
		},
	)
	var wg sync.WaitGroup

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

			// Assert:
			if testCase.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, testCase.expectedErr.Error())
			}
			require.ElementsMatch(
				t,
				testCase.expected,
				mockStream.SentEvents,
			)
		}(testCase)
	}
	wg.Wait()
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
	t.Parallel()
	// Arrange:
	currentTime := time.Now().UnixMilli()

	type testStruct struct {
		name     string
		event    models.Event
		expected *eventcrtl.Event
	}

	var testTable []testStruct
	binNilUID, _ := uuid.Nil.MarshalBinary()
	testTable = append(testTable, testStruct{
		name:     "Test 0",
		event:    models.Event{},
		expected: &eventcrtl.Event{EventId: binNilUID},
	})
	for i := 1; i < 100; i++ {
		eId := uuid.New()
		binUID, _ := eId.MarshalBinary()
		testTable = append(testTable,
			testStruct{
				name: fmt.Sprintf("Test %d", i),
				event: models.Event{
					SenderID: int64(i),
					ID:       eId,
					Time:     currentTime,
					Name:     fmt.Sprintf("User %d", i),
				},
				expected: &eventcrtl.Event{
					SenderId: int64(i),
					EventId:  binUID,
					Time:     currentTime,
					Name:     fmt.Sprintf("User %d", i),
				},
			},
		)
	}
	var wg sync.WaitGroup

	// Act:
	for _, testCase := range testTable {
		wg.Add(1)
		go func(testCase testStruct) {
			defer wg.Done()
			actual := service.AccumulateEvent(testCase.event)
			// Assert:
			require.Equal(t, testCase.expected, actual)
		}(testCase)
	}
	wg.Wait()
}

func TestPassedEvents(t *testing.T) {
	t.Parallel()
	// Arrange:
	currentTime := time.Now().UTC()

	type testStruct struct {
		name     string
		argtime  time.Time
		req      []*eventcrtl.MakeEventRequest
		expected int
	}

	var testTable []testStruct
	var requests []*eventcrtl.MakeEventRequest

	for i := 1; i < 10; i++ {
		requests = append(requests,
			&eventcrtl.MakeEventRequest{
				SenderId: int64(i),
				Time:     currentTime.AddDate(0, i, 0).UnixMilli(),
				Name:     fmt.Sprintf("User %d", i),
			},
		)
	}
	testTable = append(testTable,
		testStruct{
			name:     "Test1",
			argtime:  currentTime,
			req:      requests,
			expected: 9,
		},
		testStruct{
			name:     "Test2",
			argtime:  currentTime.AddDate(0, 5, 0),
			req:      requests,
			expected: 4,
		},
		testStruct{
			name:     "Test3",
			argtime:  currentTime.AddDate(5, 0, 0),
			req:      requests,
			expected: 0,
		},
	)
	var wg sync.WaitGroup

	// Act:
	for _, testCase := range testTable {
		wg.Add(1)
		go func(testCase testStruct) {
			defer wg.Done()
			s := service.RunEventsService(testCase.req...)
			actual := s.PassedEvents(testCase.argtime)
			// Assert:
			require.Equal(t, testCase.expected, actual)
		}(testCase)
	}
	wg.Wait()
}
