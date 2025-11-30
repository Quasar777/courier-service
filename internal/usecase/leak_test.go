package usecase

import (
	"os"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	// Проверит, что после ВСЕХ тестов в пакете usecase
	// не осталось лишних горутин
	goleak.VerifyTestMain(m)
	os.Exit(m.Run())
}
