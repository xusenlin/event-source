# event-source

## Install
```
go get github.com/xusenlin/event-source@v0.0.2-alpha
```
## Usage

```go
package task

import eventSource "github.com/xusenlin/event-source"

type EventData struct {
	//...more
}

var eventInstance = eventSource.New[*EventData]()

func PublishEventMsg() {
	eventInstance.PublishMsg("message", &EventData{})
}

```

```go
package task

import (
	"github.com/gin-gonic/gin"
	"io"
	"marewood/internal/user"
	"strconv"
)

func EventSource(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	claims, err := user.JwtParseToken(c.Query("token"))

	if err != nil {
		c.SSEvent("error", "token error")
		return
	}
	strUserId := strconv.Itoa(int(claims.ID))

	eventInstance.Subscribe(strUserId)
	defer eventInstance.CancelSubscribe(strUserId)

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-eventInstance.ReceiveMsg(strUserId); ok {
			c.SSEvent(msg.Type, msg.Data)
			return true
		}
		return false
	})
}

```
