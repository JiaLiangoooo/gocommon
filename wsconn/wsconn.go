package wsconn

import (
	"github.com/JiaLiangoooo/gocommon/logger"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	//pongWait = 10 * time.Second

	// SendMessage pings to peer with this period. Must be less than pongWait.
	//pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024
)

type WSConn struct {
	wsc    *websocket.Conn // websocket conn
	WscKey string          // 特殊标识,在日志中需要用到
	//receiveC chan Receive    // 收到客户端消息的Chan
	ctrlC chan interface{} // 需要发送给前端的Channel
	ttl   time.Duration    // ttl

	isClose bool // 是否关闭状态

	notLink *time.Timer
	lock    *sync.Mutex

	// 收到断开消息(无论是Chat服务断开, 还是前端断开)
	OnDisconnect func(code int, msg string)

	// 收到前端消息类型
	OnMessage func(msgType int, msg []byte) (isPing bool)
}

// ping
type pingrecv struct{}

func newPingRecv() (ctrl *pingrecv) {
	ctrl = &pingrecv{}
	return
}

func NewWSConn(wsc *websocket.Conn, ttl time.Duration) (conn *WSConn) {
	conn = &WSConn{
		wsc:    wsc,
		WscKey: "",
		//receiveC: make(chan Receive, 10),
		ctrlC:   make(chan interface{}, 10),
		ttl:     ttl,
		notLink: time.NewTimer(ttl),
		lock:    &sync.Mutex{},
	}
	return
}

// Start 开始websocket服务
// withPing 是否发送Ping 客户端为true, 服务端为false
// 超时是否结束
func (conn *WSConn) Start(timeBreak bool) {
	logger.Infof("conn %+v start", conn)
	conn.wsc.SetCloseHandler(func(code int, text string) error {
		logger.Debugf("wsc(%s) close message(%d)(%s)", conn.WscKey, code, text)
		//conn.Close(3004, []byte{})
		return nil
	})

	//兼容标准ping
	conn.wsc.SetPingHandler(func(message string) error {
		//logger.Debugf("wsc(%s) msgType: receive ping", conn.WscKey)
		if err := conn.writeMessage(websocket.PongMessage, []byte{}); err != nil {
			logger.Error("send pong error: %+v", err)
			return err
		}
		return nil
	})

	if timeBreak {
		conn.notLink = time.NewTimer(conn.ttl)
	} else {
		conn.notLink = nil
	}
	conn.isClose = false
	go conn.readPump()
	go conn.ctrlLoop()
}

// Close 关闭Websocket, 此时只能发送一个信号给Ctrl流程,
func (conn *WSConn) Close(code int, msg string) {
	if conn.IsClose() {
		return
	}

	logger.Infof("wsc(%s) close, int(%d), msg(%s)", conn.WscKey, code, msg)
	conn.isClose = true
	err := conn.wsc.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, msg), time.Now().Add(writeWait))
	if err != nil {
		logger.Errorf("wsc(%s) close error: %v", conn.WscKey, err)
	}
}

func (conn *WSConn) Push(msg []byte) bool {
	if conn.IsClose() {
		return false
	}
	conn.ctrlC <- msg
	return true
}

func (conn *WSConn) ctrlLoop() {
	logger.Debugf("wsc(%s) ctrl loop start", conn.WscKey)
	for {
		select {
		case ctrl, ok := <-conn.ctrlC:
			if !ok {
				break
			}
			conn.procCtrl(ctrl)
		case <-conn.notLink.C:
			logger.Infof("wsc(%s) ping timeout", conn.WscKey)
			conn.Close(PING_TIMEOUT, PING_TIMEOUT_TEXT)
			break
		}
	}
}

//  procCtrl 处理消息
// 1 MQ收到的消息
// 2 ping相关消息
func (conn *WSConn) procCtrl(ctrl interface{}) (keepRun bool) {
	switch ctrl.(type) {
	case *pingrecv:
		conn.notLink.Reset(conn.ttl)
	case []byte:
		conn.writeMessage(websocket.BinaryMessage, ctrl.([]byte))
	}
	return true
}

func (conn *WSConn) writeMessage(messageType int, data []byte) error {
	conn.lock.Lock()
	//logger.Debugf("wsc(%s) msgType: %d send", conn.WscKey, messageType)
	err := conn.wsc.WriteMessage(messageType, data)
	conn.lock.Unlock()
	return err
}

// 读取消息
func (conn *WSConn) readPump() {
	defer func() {
		conn.wsc.Close()
		conn.notLink.Stop()
	}()
	conn.wsc.SetReadLimit(maxMessageSize)
	//if e := conn.wsc.SetReadDeadline(time.Now().Add(pongWait)); e != nil {
	//	logger.Errorf("wsc() set read dead line error", conn.WscKey)
	//}

	//conn.wsc.SetPongHandler(func(appData string) error {
	//	logger.Debugf("wsc(%s) receive pong (%s)", conn.WscKey, appData)
	//	if e := conn.wsc.SetReadDeadline(time.Now().Add(pongWait)); e != nil {
	//		logger.Errorf("wsc() set read dead line error", conn.WscKey)
	//	}
	//	return nil
	//})

	for {
		msgType, message, err := conn.wsc.ReadMessage()
		if err != nil {
			conn.isClose = true
			if e, ok := err.(*websocket.CloseError); ok {
				logger.Errorf("wsc() error: %v", conn.WscKey, err)
				if conn.OnDisconnect != nil {
					conn.OnDisconnect(e.Code, e.Text)
				}
			} else {
				if conn.OnDisconnect != nil {
					conn.OnDisconnect(UNEXPECTED_CLOSE_ERRORR, err.Error())
				}
			}
			break
		}
		if conn.OnMessage == nil {
		} else if conn.OnMessage(msgType, message) {
			if err := conn.writeMessage(websocket.PongMessage, []byte{}); err != nil {
				logger.Errorf("wsc() write error: %v", conn.WscKey, err)
			}
			conn.ctrlC <- newPingRecv()
		}
	}
}

// IsClose 当前是否是关闭状态 true: 关闭 false:开启
func (conn *WSConn) IsClose() bool {
	return conn.isClose
}
