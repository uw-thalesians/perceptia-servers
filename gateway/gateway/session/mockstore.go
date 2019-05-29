package session

import (
	"fmt"
	"testing"
)

// MethodName is used to reference methods exported by MockStore.
type MethodName string

// String allows a MethodName to be printed
func (m MethodName) String() string {
	return string(m)
}

// These MethodName constants provide easy access to the methods exported by MockStore to use when
// building the map of functions to pass to the MockStore.
const (
	FNSave   MethodName = "Save"
	FNGet    MethodName = "Get"
	FNDelete MethodName = "Delete"
)

// MockStore represents a sessions.Store to be used in testing functions that rely on a session Store.
//
// If a function is called that was not provided to the MockStore constructor, or added using the AddFunction
// method, a testing.T.Fatal() will be called.
type MockStore struct {
	t                  *testing.T
	testingErrorPrefix string
	fnSave             func(SessionID, interface{}) error
	fnGet              func(SessionID, interface{}) error
	fnDelete           func(SessionID) error
}

// NewMockStore constructs a new MockStore.
//
// 'tErrPrefix' should be a string that notes any specifics about the particular instance
// of MockStore. It will be prefixed to the following ": MockStore: error: ..." if an error occurs.
// Use AddFunctions to add test methods to the mock store.
func NewMockStore(t *testing.T, tErrPrefix string) *MockStore {
	return &MockStore{t: t, testingErrorPrefix: tErrPrefix}
}

// AddFunctions takes a map of functions to add to the MockStore. These functions must meet
// 'funcs' should be a map containing the functions that are expected to be called by the code
// being tested. Additionally, if any of the provided functions do not match the expected signature of
// their MockFunction counterpart a testing.T.Fatal() will be called.
//
// If any of the provided functions already exist in the MockStore they will be replaced.
func (ms *MockStore) AddFunctions(funcs map[MethodName]interface{}) {
	ms.addFunctions(funcs)
}

// UpdateErrorPrefix will replace the current tErrPrefix.
func (ms *MockStore) UpdateErrorPrefix(tErrPrefix string) {
	ms.testingErrorPrefix = tErrPrefix
}

// Save calls the mock Save function, if this function was not mocked will cause the current test to
// log an error and fail.
func (ms *MockStore) Save(sid SessionID, sessionState interface{}) error {
	if ms.fnSave == nil {
		ms.testingError("the function (Save) was not mocked")
	}
	return ms.fnSave(sid, sessionState)
}

// Get calls the mock Get function, if this function was not mocked will cause the current test to
// log an error and fail.
func (ms *MockStore) Get(sid SessionID, sessionState interface{}) error {
	if ms.fnGet == nil {
		ms.testingError("the function (Get) was not mocked")
	}
	return ms.fnGet(sid, sessionState)
}

// Delete calls the mock Delete function, if this function was not mocked will cause the current test to
// log an error and fail.
func (ms *MockStore) Delete(sid SessionID) error {
	if ms.fnDelete == nil {
		ms.testingError("the function (Delete) was not mocked")
	}
	return ms.fnDelete(sid)
}

// addFunctions will take a map of functions and add them to this MockStore. If a provided function
// does not meet the required signature for that function a *testing.T.Fatal() will be called.
func (ms *MockStore) addFunctions(funcs map[MethodName]interface{}) *MockStore {
	if len(funcs) == 0 {
		ms.testingError("no functions were provided in 'funcs', please provide at least one function to mock")
	}
	for fnName, fn := range funcs {
		switch fnName {
		case FNSave:
			fnAdd, ok := fn.(func(SessionID, interface{}) error)
			if !ok {
				ms.testingError(fmt.Sprintf("the function supplied for %s must match: 'func(SessionID, "+
					"interface{}) error'", FNSave))
			}
			ms.fnSave = fnAdd
		case FNGet:
			fnAdd, ok := fn.(func(SessionID, interface{}) error)
			if !ok {
				ms.testingError(fmt.Sprintf("the function supplied for %s must match: 'func(SessionID, "+
					"interface{}) error'", FNGet))
			}
			ms.fnGet = fnAdd
		case FNDelete:
			fnAdd, ok := fn.(func(SessionID) error)
			if !ok {
				ms.testingError(fmt.Sprintf("the function supplied for %s must match: 'func(SessionID, "+
					"interface{}) error'", FNDelete))
			}
			ms.fnDelete = fnAdd
		default:
			ms.testingError(fmt.Sprintf("the function name (%s) is not a function of a sessions."+
				"MockStore'", fnName))
		}
	}
	return ms
}

// testingError will call testing.T.Fatal() if an error occurs with the testing code.
// If testingErrorPrefix was provided, it will be prefixed to the error message.
func (ms *MockStore) testingError(message string) {
	ms.t.Fatalf("MockStore: error: %s", message)
}
