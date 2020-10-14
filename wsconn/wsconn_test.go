package wsconn

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestWSConnStartWithoutPing(t *testing.T) {
	var mocTTL = 3 * time.Second
	http.HandleFunc("/moc_ws_server", func(writer http.ResponseWriter, request *http.Request) {
		wsc, err := websocket.Upgrade(writer, request, nil, 1024, 1024)
		assert.Nil(t, err)
		wsconn := NewWSConn(wsc, mocTTL)
		wsconn.OnDisconnect = func(code int, msg string) {
			assert.Equal(t, PING_TIMEOUT, code)
		}
		wsconn.OnMessage = nil
		wsconn.Start(true)
	})
	var mocCloseC chan interface{} = make(chan interface{})
	mockServer := &http.Server{
		Addr:    ":23233",
		Handler: nil,
	}
	go func() {
		err := mockServer.ListenAndServe()
		mocCloseC <- err
	}()
	time.Sleep(time.Second * 1)
	client := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	_, _, err := client.Dial("ws://127.0.0.1:23233/moc_ws_server", nil)
	assert.Nil(t, err)
	time.Sleep(time.Second * 3)
}

func TestConnStartWithPing(t *testing.T) {
	var mocTTL = 5 * time.Second
	http.HandleFunc("/moc_ws_server_with_ping", func(writer http.ResponseWriter, request *http.Request) {
		wsc, err := websocket.Upgrade(writer, request, nil, 1024, 1024)
		assert.Nil(t, err)
		wsconn := NewWSConn(wsc, mocTTL)
		wsconn.OnDisconnect = func(code int, msg string) {
			assert.Equal(t, PING_TIMEOUT, code)
			assert.True(t, wsconn.IsClose())
		}
		wsconn.OnMessage = nil
		wsconn.Start(true)
		assert.False(t, wsconn.isClose)
	})
	var mocCloseC chan interface{} = make(chan interface{})
	mockServer := &http.Server{
		Addr:    ":23233",
		Handler: nil,
	}
	go func() {
		err := mockServer.ListenAndServe()
		mocCloseC <- err
	}()
	time.Sleep(time.Second * 1)
	client := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	dial, _, err := client.Dial("ws://127.0.0.1:23233/moc_ws_server_with_ping", nil)
	assert.Nil(t, err)
	dial.WriteMessage(websocket.PingMessage, []byte{})
	time.Sleep(time.Second * 3)
}

func TestConnStartWithClientStop(t *testing.T) {
	var mocTTL = 5 * time.Second
	http.HandleFunc("/moc_ws_server", func(writer http.ResponseWriter, request *http.Request) {
		wsc, err := websocket.Upgrade(writer, request, nil, 1024, 1024)
		assert.Nil(t, err)
		wsconn := NewWSConn(wsc, mocTTL)
		wsconn.OnDisconnect = func(code int, msg string) {
			assert.Equal(t, websocket.CloseNoStatusReceived, code)
			assert.True(t, wsconn.IsClose())
		}
		wsconn.OnMessage = nil
		wsconn.Start(true)
		assert.False(t, wsconn.isClose)
	})
	var mocCloseC chan interface{} = make(chan interface{})
	mockServer := &http.Server{
		Addr:    ":23233",
		Handler: nil,
	}
	go func() {
		err := mockServer.ListenAndServe()
		mocCloseC <- err
	}()
	time.Sleep(time.Second * 1)
	client := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	dial, _, err := client.Dial("ws://127.0.0.1:23233/moc_ws_server", nil)
	assert.Nil(t, err)
	dial.WriteMessage(websocket.PingMessage, []byte{})
	dial.WriteMessage(websocket.CloseMessage, []byte{})
	time.Sleep(time.Second * 3)
}
