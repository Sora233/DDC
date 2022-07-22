package middleware

import (
	"github.com/Sora233/DDC/server/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func WrapperMiddleware(component string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := GetRequestID(c)
		UpdateLog(c, logrus.WithFields(logrus.Fields{
			"component":  component,
			"request_id": requestID,
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"client_ip":  c.Request.RemoteAddr,
		}))
		start := time.Now()

		c.Next()

		iresp, ok := c.Get(KeyResp)
		if !ok {
			return
		}

		gei := iresp.(model.WithGenericErrInfo).GetGenericErrInfo()
		gei.RequestID = requestID

		c.JSON(model.HTTPCode(gei.ErrorCode), iresp)

		logBuf := GetLog(c).WithFields(logrus.Fields{
			"latency":     time.Since(start),
			"body_size":   c.Writer.Size(),
			"error_code":  gei.ErrorCode,
			"status_code": c.Writer.Status(),
		})
		if gei.ErrorCode == model.DDCOK {
			logBuf.Info("OK")
		} else {
			logBuf.Error(gei.ErrorMsg)
		}
	}
}
