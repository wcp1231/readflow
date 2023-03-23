package middleware

import (
	"fmt"
	"github.com/ncarlier/readflow/pkg/config"

	"github.com/rs/zerolog/log"
)

const tpl = "using %s as authentication backend"

// Auth is a middleware to authenticate HTTP request
func Auth(config config.GlobalConfig) Middleware {
	method := config.AuthN
	switch method {
	case "mock":
		log.Info().Msg(fmt.Sprintf(tpl, "Mock"))
		return MockAuth
	case "proxy":
		log.Info().Msg(fmt.Sprintf(tpl, "Proxy"))
		return ProxyAuth
	case "basic":
		log.Info().Msg(fmt.Sprintf(tpl, "Basic"))
		return NewBasicAuth(config.BasicAuthUser, config.BasicAuthPass)
	default:
		log.Info().Str("authority", method).Msg(fmt.Sprintf(tpl, "OpenID Connect"))
		return OpenIDConnectJWTAuth(method)
	}
}
