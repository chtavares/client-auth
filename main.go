package main

import (
	"context"
	"net/http"

	"github.com/ory/hydra/sdk/go/hydra"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
)

func HandlerSlash(c echo.Context) error {
	config := &oauth2.Config{
		ClientID:     "claudio",
		ClientSecret: "auth-secret",
		RedirectURL:  "http://localhost:1323/callback",
		Scopes:       []string{"hydra.clients"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:9000/oauth2/auth",
			TokenURL: "http://localhost:9000/oauth2/token",
		},
	}

	authURL := config.AuthCodeURL("0807edf7d85e5d")

	return c.Redirect(http.StatusMovedPermanently, authURL)

}

func HandlerCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, "NO code")
	}

	client, err := hydra.NewSDK(&hydra.Configuration{
		ClientID:     "auth-server",
		ClientSecret: "auth-secret",
		PublicURL:    "http://localhost:9000",
		AdminURL:     "http://localhost:9001",
		Scopes:       []string{"hydra.clients"},
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := context.Background()

	config := &oauth2.Config{
		ClientID:     "claudio",
		ClientSecret: "auth-secret",
		RedirectURL:  "http://localhost:1323/callback",
		Scopes:       []string{"hydra.clients"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:9000/oauth2/auth",
			TokenURL: "http://localhost:9000/oauth2/token",
		},
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	introspect, _, err := client.IntrospectOAuth2Token(token.AccessToken, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, introspect.Ext)
}

func main() {
	e := echo.New()

	e.GET("/", HandlerSlash)
	e.GET("/callback", HandlerCallback)

	e.Start(":1323")

}
