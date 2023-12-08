# event-source

## Install
```
go get github.com/xusenlin/event-source@v0.0.2-alpha
```
## Usage

```go
package event

import eventSource "github.com/xusenlin/event-source"

type TaskData struct {
	
}
var Task = eventSource.New[*TaskData]()

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
	strUserId := strconv.Itoa(int(claims.ID))

	event.Task.Subscribe(strUserId)
	defer event.Task.CancelSubscribe(strUserId)

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-event.Task.ReceiveMsg(strUserId); ok {
			c.SSEvent(msg.Type, msg.Data)
			return true
		}
		return false
	})
}
```

publishMs to All User
```go
Task.PublishMsg("message",&TaskData{})
```