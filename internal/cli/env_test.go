package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func clearEnv(t *testing.T) {
	t.Helper()
	vars := []string{
		"KEY1", "KEY2", "KEY3",
		"PGXUSER", "PGUSER",
		"PGXPASSWORD", "PGPASSWORD",
		"PGXHOST", "PGHOST",
		"PGXDATABASE", "PGDATABASE",
		"PGXPORT", "PGPORT",
	}
	for _, v := range vars {
		t.Setenv(v, "")
	}
}

func Test_firstEnv(t *testing.T) {
	tests := []struct {
		name     string
		keys     []string
		env      map[string]string
		expected string
	}{
		{
			name:     "first key exists",
			keys:     []string{"KEY1", "KEY2"},
			env:      map[string]string{"KEY1": "val1", "KEY2": "val2"},
			expected: "val1",
		},
		{
			name:     "second key exists",
			keys:     []string{"KEY1", "KEY2"},
			env:      map[string]string{"KEY2": "val2"},
			expected: "val2",
		},
		{
			name:     "no keys exist",
			keys:     []string{"KEY1", "KEY2"},
			env:      map[string]string{"KEY3": "val3"},
			expected: "",
		},
		{
			name:     "empty keys",
			keys:     []string{},
			env:      map[string]string{"KEY1": "val1"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			assert.Equal(t, tt.expected, firstEnv(tt.keys...))
		})
	}
}

func Test_getUserFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "PGXUSER takes precedence",
			env:      map[string]string{"PGXUSER": "user1", "PGUSER": "user2"},
			expected: "user1",
		},
		{
			name:     "PGUSER used if PGXUSER is empty",
			env:      map[string]string{"PGUSER": "user2"},
			expected: "user2",
		},
		{
			name:     "nothing set",
			env:      map[string]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			assert.Equal(t, tt.expected, getUserFromEnv())
		})
	}
}

func Test_getPasswordFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "PGXPASSWORD takes precedence",
			env:      map[string]string{"PGXPASSWORD": "pass1", "PGPASSWORD": "pass2"},
			expected: "pass1",
		},
		{
			name:     "PGPASSWORD used if PGXPASSWORD is empty",
			env:      map[string]string{"PGPASSWORD": "pass2"},
			expected: "pass2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			assert.Equal(t, tt.expected, getPasswordFromEnv())
		})
	}
}

func Test_getHostFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "PGXHOST takes precedence",
			env:      map[string]string{"PGXHOST": "host1", "PGHOST": "host2"},
			expected: "host1",
		},
		{
			name:     "PGHOST used if PGXHOST is empty",
			env:      map[string]string{"PGHOST": "host2"},
			expected: "host2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			assert.Equal(t, tt.expected, getHostFromEnv())
		})
	}
}

func Test_getDatabaseFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "PGXDATABASE takes precedence",
			env:      map[string]string{"PGXDATABASE": "db1", "PGDATABASE": "db2"},
			expected: "db1",
		},
		{
			name:     "PGDATABASE used if PGXDATABASE is empty",
			env:      map[string]string{"PGDATABASE": "db2"},
			expected: "db2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			assert.Equal(t, tt.expected, getDatabaseFromEnv())
		})
	}
}

func Test_getPortFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected uint16
	}{
		{
			name:     "PGXPORT takes precedence",
			env:      map[string]string{"PGXPORT": "5432", "PGPORT": "5433"},
			expected: 5432,
		},
		{
			name:     "PGPORT used if PGXPORT is empty",
			env:      map[string]string{"PGPORT": "5433"},
			expected: 5433,
		},
		{
			name:     "invalid port returns 0",
			env:      map[string]string{"PGXPORT": "abc"},
			expected: 0,
		},
		{
			name:     "port out of range returns 0",
			env:      map[string]string{"PGXPORT": "70000"},
			expected: 0,
		},
		{
			name:     "nothing set returns 0",
			env:      map[string]string{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			assert.Equal(t, tt.expected, getPortFromEnv())
		})
	}
}
