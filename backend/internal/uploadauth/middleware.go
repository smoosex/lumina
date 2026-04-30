package uploadauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smoosex/lumina/backend/internal/config"
	"go.uber.org/fx"
)

var Module = fx.Module("uploadauth", fx.Provide(New))

type Middleware struct {
	cfg        config.UploadConfig
	allowedIPs map[string]struct{}
}

func New(cfg config.Config) *Middleware {
	allowedIPs := make(map[string]struct{}, len(cfg.Upload.AllowedIPs))
	for _, ip := range cfg.Upload.AllowedIPs {
		allowedIPs[ip] = struct{}{}
	}

	return &Middleware{
		cfg:        cfg.Upload,
		allowedIPs: allowedIPs,
	}
}

func (m *Middleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(m.allowedIPs) > 0 {
			host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
			if err != nil {
				host = c.ClientIP()
			}
			if _, ok := m.allowedIPs[host]; !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "upload source is not allowed",
				})
				return
			}
		}

		body, err := io.ReadAll(http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			m.cfg.MaxMarkdownBytes,
		))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "request body is too large",
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		if !m.valid(c.Request, body) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid upload signature",
			})
			return
		}

		c.Next()
	}
}

func (m *Middleware) valid(req *http.Request, body []byte) bool {
	timestamp := req.Header.Get("X-Lumina-Timestamp")
	bodyHash := req.Header.Get("X-Lumina-Body-SHA256")
	signature := req.Header.Get("X-Lumina-Signature")
	if timestamp == "" || bodyHash == "" || signature == "" {
		return false
	}

	unix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	requestTime := time.Unix(unix, 0)
	if time.Since(requestTime).Abs() > m.cfg.TimestampSkew {
		return false
	}

	sum := sha256.Sum256(body)
	expectedBodyHash := hex.EncodeToString(sum[:])
	if subtle.ConstantTimeCompare(
		[]byte(bodyHash),
		[]byte(expectedBodyHash),
	) != 1 {
		return false
	}

	payload := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		req.Method,
		req.URL.Path,
		timestamp,
		bodyHash,
	)
	mac := hmac.New(sha256.New, []byte(m.cfg.SignatureSecret))
	_, _ = mac.Write([]byte(payload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return subtle.ConstantTimeCompare(
		[]byte(signature),
		[]byte(expectedSignature),
	) == 1
}
