package sanitizer_test

import (
	"fmt"
	"testing"

	"github.com/satriahrh/animated-waffle/sanitizer"
	"github.com/stretchr/testify/require"
)

func TestSanitizeFetch(t *testing.T) {
	for _, row := range []struct {
		Name     string
		Input    []string
		Expected error
	}{
		{"Empty", []string{}, fmt.Errorf("no url given")},
		{"OneInvalid", []string{"not-a-file"}, fmt.Errorf("invalid url")},
		{"OneInvalidHostOnly", []string{"google.com"}, fmt.Errorf("invalid url")},
		{"OneSuccessSample", []string{"https://google.com"}, nil},
		{"OneSuccessWithPath", []string{"https://google.com/webpage"}, nil},
		{"MultipleHaveInvalid", []string{"https://google.com/", "invalid"}, fmt.Errorf("invalid url")},
		{"MultipleSuccess", []string{"https://google.com/", "https://facebook.com"}, nil},
	} {
		t.Run(row.Name, func(t *testing.T) {
			actual := sanitizer.SanitizeFetch(nil, row.Input)
			require.Equal(t, row.Expected, actual)
		})
	}
}
