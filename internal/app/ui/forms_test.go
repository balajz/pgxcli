package ui

import (
	"testing"

	"charm.land/huh/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModel_DefaultPort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		wantPort string
	}{
		{name: "empty port defaults", input: "", wantPort: "5432"},
		{name: "zero port defaults", input: "0", wantPort: "5432"},
		{name: "spaced zero defaults", input: " 0 ", wantPort: "5432"},
		{name: "keeps provided port", input: "6543", wantPort: "6543"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := newModel("postgres", "user", "localhost", tc.input)
			require.NotNil(t, m)
			assert.Equal(t, tc.wantPort, m.values.Port)
			require.NotNil(t, m.form)
		})
	}
}

func TestModelResult(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		state     huh.FormState
		port      string
		wantErr   string
		wantValue ConnectionValues
	}{
		{
			name:    "aborted when form incomplete",
			state:   huh.StateNormal,
			port:    "5432",
			wantErr: "connection form aborted",
		},
		{
			name:    "invalid port on completed form",
			state:   huh.StateCompleted,
			port:    "70000",
			wantErr: "invalid port",
		},
		{
			name:  "returns values when completed",
			state: huh.StateCompleted,
			port:  "5432",
			wantValue: ConnectionValues{
				Database: "postgres",
				Username: "user",
				Host:     "localhost",
				Port:     "5432",
				Password: "secret",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := newModel("postgres", "user", "localhost", tc.port)
			m.values.Password = "secret"
			m.form.State = tc.state

			got, err := m.result()
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantValue, got)
		})
	}
}

func TestFallback(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty to none", in: "", want: "(None)"},
		{name: "spaces to none", in: "   ", want: "(None)"},
		{name: "trimmed value", in: "  postgres  ", want: "postgres"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, fallback(tc.in))
		})
	}
}

func TestMaskedPassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty to none", in: "", want: "(None)"},
		{name: "single char", in: "a", want: "•"},
		{name: "three chars", in: "abc", want: "•••"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, maskedPassword(tc.in))
		})
	}
}

func TestValidatePort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		in      string
		wantErr string
	}{
		{name: "empty", in: "", wantErr: "port is required"},
		{name: "not number", in: "abc", wantErr: "port must be a number"},
		{name: "too low", in: "0", wantErr: "port must be between 1 and 65535"},
		{name: "too high", in: "65536", wantErr: "port must be between 1 and 65535"},
		{name: "valid", in: "5432"},
		{name: "valid with spaces", in: " 5432 "},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := validatePort(tc.in)
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
