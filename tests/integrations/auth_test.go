package integrations

import (
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests"
	"github.com/notefan-golang/models/responses"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuth(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	client := httpClient()

	password := "11112222"
	passwordBcrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.Nil(err)
	user := entities.User{
		Name:     faker.Name(),
		Email:    faker.Email(),
		Password: string(passwordBcrypt),
	}

	t.Run("Register", func(t *testing.T) {
		reqBody := requests.AuthRegister{
			Name:            user.Name,
			Email:           user.Email,
			Password:        password,
			ConfirmPassword: password,
		}
		reqBodyStr, err := helper.JSONStructToJSONStr(reqBody)
		require.Nil(err)

		expectedHttpCode := http.StatusCreated
		apitest.New().
			EnableNetworking(client).
			Post("http://localhost:8080/api/users/register").
			JSON(reqBodyStr).
			Expect(t).
			Status(expectedHttpCode).
			Assert(jsonpath.Equal("code", float64(expectedHttpCode))).
			Assert(jsonpath.Equal("status", responses.StatusSuccess)).
			Assert(jsonpath.Present("message")).
			Assert(jsonpath.Root("user").
				Present("id").
				Equal("name", reqBody.Name).
				Equal("email", reqBody.Email).
				NotPresent("password").
				End(),
			).
			End()
	})

	t.Run("Login", func(t *testing.T) {
		reqBody := requests.AuthLogin{
			Email:    user.Email,
			Password: password,
		}
		reqBodyStr, err := helper.JSONStructToJSONStr(reqBody)
		require.Nil(err)

		expectedHttpCode := http.StatusOK
		apitest.New().
			EnableNetworking(client).
			Post("http://localhost:8080/api/users/login").
			JSON(reqBodyStr).
			Expect(t).
			Status(expectedHttpCode).
			Assert(jsonpath.Equal("code", float64(expectedHttpCode))).
			Assert(jsonpath.Equal("status", responses.StatusSuccess)).
			Assert(jsonpath.Present("message")).
			Assert(jsonpath.Root("user").
				Present("id").
				Equal("name", user.Name).
				Equal("email", reqBody.Email).
				Present("access_token").
				End(),
			).
			End()
	})

	t.Run("Logout", func(t *testing.T) {
		reqBody := requests.AuthLogin{
			Email:    user.Email,
			Password: password,
		}
		reqBodyStr, err := helper.JSONStructToJSONStr(reqBody)
		require.Nil(err)

		expectedHttpCode := http.StatusOK
		apitest.New().
			EnableNetworking(client).
			Delete("http://localhost:8080/api/users/logout").
			JSON(reqBodyStr).
			Expect(t).
			Status(expectedHttpCode).
			Assert(jsonpath.Equal("code", float64(expectedHttpCode))).
			Assert(jsonpath.Equal("status", responses.StatusSuccess)).
			Assert(jsonpath.Present("message")).
			End()
	})

}
