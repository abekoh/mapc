package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func requireNoErrors(t *testing.T, errs []error) {
	t.Helper()
	for _, err := range errs {
		require.Nil(t, err)
	}
}
