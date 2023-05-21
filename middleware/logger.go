package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func AccessLog() gin.HandlerFunc {
	logger := Logger()

	return func(c *gin.Context) {
		t := time.Now()

		// 初始化bodyLogWriter
		bodyLogWriter := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = bodyLogWriter

		//获取请求信息
		requestBody := getRequestBody(c)

		c.Next()

		// 执行时间
		latencyTime := time.Since(t)

		//响应内容
		responseBody := bodyLogWriter.body.String()

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 请求参数
		//params := c.Request.

		//headers := ""
		//for k, v := range c.Request.Header {
		//	headers += k + ":" + strings.Join(v, ",") + "&"
		//}
		//
		//headers := c.Request.Header.Get("X-uid") + ":" + c.Request.Header.Get("X-token")
		//body, _ := ioutil.ReadAll(c.Request.Body)

		//日志格式
		logger.Infof("| %s | %3d | %13v | %15s | %s | %s | %s | %s |",
			c.Request.Context().Value("trace_id"),
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			requestBody,
			responseBody,
			reqUri,
		)
	}
}

func getRequestBody(ctx *gin.Context) interface{} {
	switch ctx.Request.Method {
	case http.MethodGet:
		return ctx.Request.URL.Query()
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		var bodyBytes []byte
		bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return nil
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return string(bodyBytes)
	}

	return nil
}
