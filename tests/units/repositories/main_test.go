package repositories

import (
	"testing"

	"github.com/notefan-golang/tests"
)

func TestMain(m *testing.M) {
	tests.Setup()
	defer func() {
		tests.Teardown()
	}()

	// Run tests
	m.Run()
}
