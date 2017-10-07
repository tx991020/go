package common

import (
	"crypto/rsa"

	"io/ioutil"
	"log"


	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go/request"



)

type AppClaims struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

const (
	// openssl genrsa -out app.rsa 1024
	privKeyPath = "keys/tm.rsa"
	// openssl rsa -in app.rsa -pubout > app.rsa.pub
	pubKeyPath = "keys/tm.rsa.pub"
)

// Private key for signing and public key for verification
var (
	//verifyKey, signKey []byte
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func initKeys() {

	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
	}
}


func GenerateJWT(name, role string) (string, error) {
	// Create the Claims
	claims := AppClaims{
		name,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
			Issuer:    "admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	ss, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}


func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from request
		token, err := request.ParseFromRequestWithClaims(c.Request, request.OAuth2Extractor, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

		if err != nil {
			switch err.(type) {
			case *jwt.ValidationError: // JWT validation error
				vErr := err.(*jwt.ValidationError)

				switch vErr.Errors {
				case jwt.ValidationErrorExpired: //JWT expired
					DisplayAppError(
						c,
						err,
						"Access Token is expired, get a new Token",
						401,
					)
				default:
					DisplayAppError(
						c,
						err,
						"Error while parsing the Access Token!",
						500,
					)
				}
			default:
				DisplayAppError(
					c,
					err,
					"Error while parsing Access Token!",
					500)
			}
			c.Abort()
		}
		if token.Valid {
			c.Set("user", token.Claims.(*AppClaims).UserName)
			c.Next()
		} else {
			DisplayAppError(
				c,
				err,
				"Invalid Access Token",
				401,
			)
			c.Abort()
		}
	}
}
