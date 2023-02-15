package integrations

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/notefan-golang/config"
	"github.com/notefan-golang/containers"

	"github.com/stretchr/testify/assert"
)

func TestAuthRegisterLoginLogout(t *testing.T) {
	assert := assert.New(t)

	app := config.InitializeTestApp()
	authController := containers.InitializeAuthController(app.DB)

	recorder := httptest.NewRecorder()

	t.Run("Register", func(t *testing.T) {
		registerRequest, err := http.NewRequest(http.MethodPost, "/api/users/register", nil)
		assert.Nil(err)

		authController.Register(recorder, registerRequest)
		registerResponse := recorder.Result()
		defer registerResponse.Body.Close()

		body, err := io.ReadAll(registerResponse.Body)
		assert.Nil(err)

		t.Error(body)
	})

	assert.True(true)
}
