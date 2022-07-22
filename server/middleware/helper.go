package middleware

import (
	"context"
	"github.com/Sora233/DDC/server/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetLog(c context.Context) logrus.FieldLogger {
	if c != nil {
		l := c.Value(KeyLog)
		if l != nil {
			return l.(logrus.FieldLogger)
		}
	}
	return logrus.StandardLogger()
}

func UpdateLog(c *gin.Context, entry logrus.FieldLogger) {
	if c == nil {
		return
	}
	c.Set(KeyLog, entry)
}

func SetResp(c *gin.Context, resp model.WithGenericErrInfo) {
	if c == nil {
		return
	}
	c.Set(KeyResp, resp)
}

func GetRequestID(c context.Context) string {
	if c == nil {
		return ""
	}
	return c.Value(KeyRequestID).(string)
}
