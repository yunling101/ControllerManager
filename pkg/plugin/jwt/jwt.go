package jwt

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/toolkits/file"
	"os"
)

type parsed struct {
	secretKey string
}

func New(secretKey string) *parsed {
	return &parsed{secretKey: secretKey}
}

func Env() *parsed {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		panic(any("SECRET_KEY environment variable not specified"))
	}
	return &parsed{secretKey: secretKey}
}

func File(filename string) *parsed {
	if !file.IsExist(filename) {
		panic(any("SECRET_KEY File does not exist"))
	}
	secretKey, err := os.ReadFile(filename)
	if err != nil {
		panic(any("SECRET_KEY ReadFile error"))
	}
	return &parsed{secretKey: string(secretKey)}
}

func (c *parsed) Registered() (sign string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	sign, err = token.SignedString([]byte(c.secretKey))
	return
}

func (c *parsed) ParseSigned(sign string) error {
	parsedToken, err := jwt.Parse(sign, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.secretKey), nil
	})
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		err = fmt.Errorf("validation failed")
		return err
	}
	return nil
}

func (c *parsed) AuthToken(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	signedToken := request.HeaderParameter("Token")
	if err := c.ParseSigned(signedToken); err != nil {
		_ = response.WriteAsJson(map[string]interface{}{
			"code":    1001,
			"message": "invalid token authentication",
		})
		return
	}
	chain.ProcessFilter(request, response)
}
