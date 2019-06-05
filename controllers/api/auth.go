package api

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo"

	nonce_caches "github.com/kilfu0701/echo_app/caches/nonce"
	"github.com/kilfu0701/echo_app/core"
	c_err "github.com/kilfu0701/echo_app/core/errors"
	"github.com/kilfu0701/echo_app/models/api_users"
)

func Hello(c echo.Context) error {
	ac := c.(*core.AppContext)

	XUserIdentifier := ac.Request().Header.Get("X-User-Identifier")
	if XUserIdentifier == "" {
		return ac.JSON(http.StatusBadRequest, emptyHelloResult())
	}

	username := "user:" + XUserIdentifier

	log.Printf("searching for user : %s", username)

	// search in api_users
	au := api_users.New(*ac)
	apiUserDoc, err := au.FindByName(username)
	if err != nil {
		if c_err.IsDocumentNotFound(err) {
			// create new document...
			ts := core.Microtime()
			uid := core.Uniqid("new_user", true)
			input_str := fmt.Sprintf("secret ... %s%s%s", ts, uid, username)

			h := sha1.New()
			io.WriteString(h, input_str)
			hashed := string(h.Sum(nil))
			hashedString := fmt.Sprintf("%x", hashed)
			userSecretHash := core.UrlsafeBase64Encode(hashedString)

			if doc, err := au.CreateUser(username, userSecretHash); err != nil {
				log.Printf("failed create user. err = %++v", err)
				return ac.JSON(http.StatusBadRequest, emptyHelloResult())
			} else {
				res := map[string]interface{}{
					"userKey":        doc.UserKey,
					"userSecretHash": doc.UserSecretHash,
				}
				return ac.JSON(http.StatusOK, res)
			}
		} else {
			log.Printf("err %++v", err)
			return ac.JSON(http.StatusBadRequest, emptyHelloResult())
		}
	}

	if err := au.UpdateLastLogin(apiUserDoc.Id); err != nil {
		log.Printf("UpdateLastLogin failed. apiUserDoc.Id=%v err=%++v", apiUserDoc.Id, err)
	}

	res := map[string]interface{}{
		"userKey":        apiUserDoc.UserKey,
		"userSecretHash": apiUserDoc.UserSecretHash,
	}
	return ac.JSON(http.StatusOK, res)
}

func Auth(c echo.Context) error {
	ac := c.(*core.AppContext)

	XUser := ac.Request().Header.Get("X-User")

	au := api_users.New(*ac)
	_, err := au.FindByKey(XUser)
	if err != nil {
		return err
	}

	//hash('sha256', 'nonce' . microtime() . uniqid(), true)
	h := sha256.New()
	input_str := fmt.Sprintf("nonce%s", core.Microtime(), core.Uniqid("", false))
	h.Write([]byte(input_str))
	hashed := h.Sum(nil)
	log.Printf("%x", h.Sum(nil))

	nonce := core.UrlsafeBase64Encode(string(hashed))

	cacheKey := "nonce." + nonce

	nc := nonce_caches.New(*ac)

	// set
	nc.SetByKey(cacheKey, &nonce_caches.NonceData{Nonce: nonce})
	_, err = nc.FindByKey(cacheKey)

	res := map[string]interface{}{
		"code":    200,
		"message": "ok",
	}
	return ac.JSON(http.StatusOK, res)
}

func AuthWithToken(c echo.Context) error {
	ac := c.(*core.AppContext)
	res := map[string]interface{}{
		"code":    200,
		"message": "ok",
	}
	return ac.JSON(http.StatusOK, res)
}

func emptyHelloResult() map[string]interface{} {
	res := map[string]interface{}{
		"userKey":        nil,
		"userSecretHash": nil,
	}
	return res
}
