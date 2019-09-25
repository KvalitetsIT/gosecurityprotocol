package securityprotocol

import (
        "net/http"
)

type HttpHandler interface {

	Handle(http.ResponseWriter, *http.Request) (int, error)
}


type SessionIdHandler interface {

	GetSessionIdFromHttpRequest(*http.Request) string
	SetSessionIdOnHttpRequest(string, *http.Request)
}


type HttpHeaderSessionIdHandler struct {

	HttpHeaderName	string
}

func (handler HttpHeaderSessionIdHandler) GetSessionIdFromHttpRequest(request *http.Request) string {

	sessionId := request.Header.Get(handler.HttpHeaderName)
	if (len(sessionId) == 0) {
		sessionId = request.Cookie(handler.HttpHeaderName)
	}
	return sessionId
}

func (handler HttpHeaderSessionIdHandler) SetSessionIdOnHttpRequest(sessionId string, request *http.Request)  {

	request.Header.Add(handler.HttpHeaderName, sessionId)
}

type MatchHandler func(r *http.Request) bool 

func MatchAllHandler (r *http.Request) bool {
	return true
}
