package e2etest_test

import (
	"testing"

	"github.com/quasilyte/quasigo/internal/testutil"
)

func TestE2E(t *testing.T) {
	runner := testutil.NewRunner(t)
	runner.Run()
}
