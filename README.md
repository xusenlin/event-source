# event-source

## Install
```
go get github.com/xusenlin/event-source@v0.0.6-alpha
```
## Usage

```go
package task

import eventSource "github.com/xusenlin/event-source"

type EventData struct {
	RepositoryId uint   `json:"repositoryId"`
	UserId       uint   `json:"userId"`
	UserName     string `json:"userName"`
	Msg          string `json:"msg"`
}

var eventInstance = eventSource.New[*EventData,uint]()

func PublishEventMsg() {
	eventInstance.PublishMsg("message", &EventData{})
	eventInstance.PublishMsgExcludeID("message", &EventData{}, 1)
	eventInstance.PublishMsgExcludeIDs("message", &EventData{}, []uint{1, 2, 3})
}

```
```go

func EventSource(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	claims, err := user.JwtParseToken(c.Query("token"))

	if err != nil {
		c.SSEvent("error", "token error")
		return
	}
	
	eventInstance.Subscribe(claims.ID)
	defer eventInstance.CancelSubscribe(claims.ID)
	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-eventInstance.ReceiveMsg(claims.ID):
			c.SSEvent(msg.Type, msg.Data)
			return true
		case <-c.Writer.CloseNotify():
			return false
		}
	})
}

```
