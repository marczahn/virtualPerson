package consciousness

import (
	"os"
	"testing"

	"github.com/marczahn/person/internal/i18n"
)

func TestMain(m *testing.M) {
	if err := i18n.Load("en"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
