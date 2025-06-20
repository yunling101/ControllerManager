package oauth

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/yunling101/ControllerManager/common"
	"math/rand"
	"strings"
	"time"
)

const (
	clientSecret = "clientSecret"
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+;<:>."
)

type authToken struct {
	Code        string
	RedirectUri string
	GrantType   string
}

func New() *authToken {
	return &authToken{}
}

func (c *authToken) GetState() common.Request {
	state := c.generateState(32)
	oauthState := c.hashStateCode(state, clientSecret)

	return common.Request{
		"state":       state,
		"oauth_state": oauthState,
	}
}

func (c *authToken) SetCode(code string) *authToken {
	c.Code = code
	return c
}

func (c *authToken) GetToken() common.Request {
	m := md5.New()
	m.Write([]byte(strings.Join([]string{c.Code, fmt.Sprintf("%v", time.Now().Unix()+7200)}, ":")))
	token := hex.EncodeToString(m.Sum(nil))

	return common.Request{
		"access_token": token, "token_type": "bearer", "scope": "",
	}
}

func (c *authToken) hashStateCode(code, clientSecret string) string {
	hash := sha256.New()
	hash.Write([]byte(code + common.Config().Yaml.OauthEncryptKey + clientSecret))
	return hex.EncodeToString(hash.Sum(nil))
}

func (c *authToken) generateState(n int) string {
	const (
		letterIdxBits = 6
		letterIdxMask = 1<<letterIdxBits - 1
		letterIdxMax  = 63 / letterIdxBits
	)

	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return base64.StdEncoding.EncodeToString(b)
}
