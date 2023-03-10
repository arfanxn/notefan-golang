package integrations

import (
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/notefan-golang/helpers/jsonh"
	"github.com/notefan-golang/models/entities"
	authReqs "github.com/notefan-golang/models/requests/auth_reqs"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/tests"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuth(t *testing.T) {
	require := require.New(t)

	password := faker.Password()
	passwordBcrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.Nil(err)
	user := entities.User{
		Name:     faker.Name(),
		Email:    faker.Email(),
		Password: string(passwordBcrypt),
	}

	t.Run("Register", func(t *testing.T) {
		reqBody := authReqs.Register{
			Name:            user.Name,
			Email:           user.Email,
			Password:        password,
			ConfirmPassword: password,
		}
		reqBodyStr, err := jsonh.ToJsonStr(reqBody)
		require.Nil(err)

		expectedHttpCode := http.StatusCreated
		apitest.New().
			EnableNetworking(tests.GetHTTPClient()).
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
		reqBody := authReqs.Login{
			Email:    user.Email,
			Password: password,
		}
		reqBodyStr, err := jsonh.ToJsonStr(reqBody)
		require.Nil(err)

		expectedHttpCode := http.StatusOK
		apitest.New().
			EnableNetworking(tests.GetHTTPClient()).
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
		reqBody := authReqs.Login{
			Email:    user.Email,
			Password: password,
		}
		reqBodyStr, err := jsonh.ToJsonStr(reqBody)
		require.Nil(err)

		expectedHttpCode := http.StatusOK
		apitest.New().
			EnableNetworking(tests.GetHTTPClient()).
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
