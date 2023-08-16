package sanitizer_test

import (
	"testing"

	"github.com/satriahrh/autify-tht/sanitizer"
	"github.com/stretchr/testify/require"
)

func TestSanitizeFetch(t *testing.T) {
	for _, row := range []struct {
		Name     string
		Input    []string
		Expected error
	}{} {
		t.Run(row.Name, func(t *testing.T) {
			actual := sanitizer.SanitizeFetch(row.Input)
			require.Equal(t, row.Expected, actual)
		})
	}
}
