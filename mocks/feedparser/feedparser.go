// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package parser_mock

import (
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"sync"
)

// Ensure, that ParserMock does implement feedparser.Parser.
// If this is not the case, regenerate this file with moq.
var _ feedparser.Parser = &ParserMock{}

// ParserMock is a mock implementation of feedparser.Parser.
//
//	func TestSomethingThatUsesParser(t *testing.T) {
//
//		// make and configure a mocked feedparser.Parser
//		mockedParser := &ParserMock{
//			ParseFunc: func(channelID string) (*feedparser.Channel, error) {
//				panic("mock out the Parse method")
//			},
//			ParseThreadSafeFunc: func(channelID string, channelChan chan *feedparser.Channel, errChan chan error, mu *sync.Mutex, wg *sync.WaitGroup)  {
//				panic("mock out the ParseThreadSafe method")
//			},
//		}
//
//		// use mockedParser in code that requires feedparser.Parser
//		// and then make assertions.
//
//	}
type ParserMock struct {
	// ParseFunc mocks the Parse method.
	ParseFunc func(channelID string) (*feedparser.Channel, error)

	// ParseThreadSafeFunc mocks the ParseThreadSafe method.
	ParseThreadSafeFunc func(channelID string, channelChan chan *feedparser.Channel, errChan chan error, mu *sync.Mutex, wg *sync.WaitGroup)

	// calls tracks calls to the methods.
	calls struct {
		// Parse holds details about calls to the Parse method.
		Parse []struct {
			// ChannelID is the channelID argument value.
			ChannelID string
		}
		// ParseThreadSafe holds details about calls to the ParseThreadSafe method.
		ParseThreadSafe []struct {
			// ChannelID is the channelID argument value.
			ChannelID string
			// ChannelChan is the channelChan argument value.
			ChannelChan chan *feedparser.Channel
			// ErrChan is the errChan argument value.
			ErrChan chan error
			// Mu is the mu argument value.
			Mu *sync.Mutex
			// Wg is the wg argument value.
			Wg *sync.WaitGroup
		}
	}
	lockParse           sync.RWMutex
	lockParseThreadSafe sync.RWMutex
}

// Parse calls ParseFunc.
func (mock *ParserMock) Parse(channelID string) (*feedparser.Channel, error) {
	if mock.ParseFunc == nil {
		panic("ParserMock.ParseFunc: method is nil but Parser.Parse was just called")
	}
	callInfo := struct {
		ChannelID string
	}{
		ChannelID: channelID,
	}
	mock.lockParse.Lock()
	mock.calls.Parse = append(mock.calls.Parse, callInfo)
	mock.lockParse.Unlock()
	return mock.ParseFunc(channelID)
}

// ParseCalls gets all the calls that were made to Parse.
// Check the length with:
//
//	len(mockedParser.ParseCalls())
func (mock *ParserMock) ParseCalls() []struct {
	ChannelID string
} {
	var calls []struct {
		ChannelID string
	}
	mock.lockParse.RLock()
	calls = mock.calls.Parse
	mock.lockParse.RUnlock()
	return calls
}

// ParseThreadSafe calls ParseThreadSafeFunc.
func (mock *ParserMock) ParseThreadSafe(channelID string, channelChan chan *feedparser.Channel, errChan chan error, mu *sync.Mutex, wg *sync.WaitGroup) {
	if mock.ParseThreadSafeFunc == nil {
		panic("ParserMock.ParseThreadSafeFunc: method is nil but Parser.ParseThreadSafe was just called")
	}
	callInfo := struct {
		ChannelID   string
		ChannelChan chan *feedparser.Channel
		ErrChan     chan error
		Mu          *sync.Mutex
		Wg          *sync.WaitGroup
	}{
		ChannelID:   channelID,
		ChannelChan: channelChan,
		ErrChan:     errChan,
		Mu:          mu,
		Wg:          wg,
	}
	mock.lockParseThreadSafe.Lock()
	mock.calls.ParseThreadSafe = append(mock.calls.ParseThreadSafe, callInfo)
	mock.lockParseThreadSafe.Unlock()
	mock.ParseThreadSafeFunc(channelID, channelChan, errChan, mu, wg)
}

// ParseThreadSafeCalls gets all the calls that were made to ParseThreadSafe.
// Check the length with:
//
//	len(mockedParser.ParseThreadSafeCalls())
func (mock *ParserMock) ParseThreadSafeCalls() []struct {
	ChannelID   string
	ChannelChan chan *feedparser.Channel
	ErrChan     chan error
	Mu          *sync.Mutex
	Wg          *sync.WaitGroup
} {
	var calls []struct {
		ChannelID   string
		ChannelChan chan *feedparser.Channel
		ErrChan     chan error
		Mu          *sync.Mutex
		Wg          *sync.WaitGroup
	}
	mock.lockParseThreadSafe.RLock()
	calls = mock.calls.ParseThreadSafe
	mock.lockParseThreadSafe.RUnlock()
	return calls
}