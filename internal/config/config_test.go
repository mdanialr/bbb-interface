package config

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeReader just fake type to satisfies io.Reader interfaces
// so it could trigger error buffer read from.
type fakeReader struct{}

func (_ fakeReader) Read(_ []byte) (_ int, _ error) {
	return 0, fmt.Errorf("this should trigger error in test")
}

func TestNewConfig(t *testing.T) {
	fakeConfigFile := `
env: prod
port: 6565
log: /home/nzk/test-app/log
`
	buf := bytes.NewBufferString(fakeConfigFile)
	t.Run("Using valid value should be pass", func(t *testing.T) {
		mod, err := NewConfig(buf)
		require.NoError(t, err)

		assert.Equal(t, "prod", mod.Env)
		assert.Equal(t, "", mod.Host)
		assert.Equal(t, uint16(6565), mod.PortNum)
		assert.Equal(t, "/home/nzk/test-app/log", mod.LogDir)
	})

	fakeConfigFile = `
env: 2134
host: 313
port: "number"
	`
	buf = bytes.NewBufferString(fakeConfigFile)
	t.Run("Using mismatch type should be error in yaml unmarshalling", func(t *testing.T) {
		_, err := NewConfig(buf)
		require.Error(t, err)
	})

	t.Run("Injecting fake reader should be error in buffer read from", func(t *testing.T) {
		_, err := NewConfig(fakeReader{})
		require.Error(t, err)
	})

	fakeConfigFile = `port: 1234`
	buf = bytes.NewBufferString(fakeConfigFile)
	t.Run("Port w 1234 should be no error", func(t *testing.T) {
		_, err := NewConfig(buf)
		require.NoError(t, err)
	})
}

func TestIsDifferentHash(t *testing.T) {
	fakeConfigOne := `this is the first file`
	fakeConfigTwo := `this is the first file`
	bufOne := bytes.NewBufferString(fakeConfigOne)
	bufTwo := bytes.NewBufferString(fakeConfigTwo)

	t.Run("Using same file should be equal", func(t *testing.T) {
		out, err := IsDifferentHash(bufOne, bufTwo)
		require.NoError(t, err)
		assert.True(t, out)
	})

	fakeConfigTwo = `this is the real second file`
	bufTwo = bytes.NewBufferString(fakeConfigTwo)
	t.Run("Using different file should not be equal", func(t *testing.T) {
		out, err := IsDifferentHash(bufOne, bufTwo)
		require.NoError(t, err)
		assert.False(t, out)
	})

	t.Run("Injecting fake reader should be error in copying first file", func(t *testing.T) {
		_, err := IsDifferentHash(fakeReader{}, bufTwo)
		require.Error(t, err)
	})

	t.Run("Injecting fake reader should be error in copying second file", func(t *testing.T) {
		_, err := IsDifferentHash(bufOne, fakeReader{})
		require.Error(t, err)
	})
}

func TestReloadConfig(t *testing.T) {
	oldMod := Model{
		Env:     "dev",
		PortNum: 1234,
		LogDir:  "/path/to/test-app/log",
	}

	// Pretend we already have one model. Then pretend that we will
	// load new one config file and repopulate old model with the
	// newly loaded. So we can compare the old model vs new model
	// to make sure that the old model successfully reloaded.

	newFakeConfigFile := `
env: prod
port: 1235
log: /var/log/webhook/log
`
	buf := bytes.NewBufferString(newFakeConfigFile)

	newMod, err := NewConfig(buf)
	require.NoError(t, err)

	t.Run("Should be no error. Then compare old v new", func(t *testing.T) {
		err := oldMod.ReloadConfig(buf)
		require.NoError(t, err)

		assert.NotEqual(t, newMod.Env, oldMod.Env)
		assert.NotEqual(t, newMod.PortNum, oldMod.PortNum)
		assert.NotEqual(t, newMod.LogDir, oldMod.LogDir)
	})

	t.Run("Injecting fake reader should be error", func(t *testing.T) {
		err := oldMod.ReloadConfig(fakeReader{})
		require.Error(t, err)
	})
}

func TestSanitization_Env(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
	}{
		{
			name:   "Env w dev should be dev",
			sample: Model{Env: "dev"},
			expect: "dev",
		},
		{
			name:   "Env w prod should be prod",
			sample: Model{Env: "prod"},
			expect: "prod",
		},
		{
			name:   "Env w/o value should be dev",
			sample: Model{},
			expect: "dev",
		},
		{
			name:   "Env w value not match either dev or prod should be dev",
			sample: Model{Env: "uu"},
			expect: "dev",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.Env)
		})
	}
}

func TestSanitization_EnvIsProd(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect bool
	}{
		{
			name:   "Env w dev should be false",
			sample: Model{Env: "dev"},
			expect: false,
		},
		{
			name:   "Env w value not match either dev or prod should be false",
			sample: Model{Env: "lol"},
			expect: false,
		},
		{
			name:   "Env w/o should be false",
			sample: Model{},
			expect: false,
		},
		{
			name:   "Env w prod should be true",
			sample: Model{Env: "prod"},
			expect: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.EnvIsProd)
		})
	}
}

func TestSanitization_Host(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
	}{
		{
			name:   "Host w localhost should be localhost",
			sample: Model{Host: "localhost"},
			expect: "localhost",
		},
		{
			name:   "Host w 120.222.46.23 should be 120.222.46.23",
			sample: Model{Host: "120.222.46.23"},
			expect: "120.222.46.23",
		},
		{
			name:   "Host w/o value should be localhost",
			sample: Model{},
			expect: "localhost",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.Host)
		})
	}
}

func TestSanitization_Port(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect uint16
	}{
		{
			name:   "Port w 2003 should be 2003",
			sample: Model{PortNum: 2003},
			expect: 2003,
		},
		{
			name:   "Port w 44444 should be 44444",
			sample: Model{PortNum: 44444},
			expect: 44444,
		},
		{
			name:   "Port w/o value should be 6767",
			sample: Model{},
			expect: 6767,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.PortNum)
		})
	}
}

func TestSanitization_RandomLen(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect uint8
	}{
		{
			name:   "Random length w 24 should be 24",
			sample: Model{RandomLen: 24},
			expect: 24,
		},
		{
			name:   "Port w 64 should be 64",
			sample: Model{RandomLen: 64},
			expect: 64,
		},
		{
			name:   "Port w/o value should be default to 8",
			sample: Model{},
			expect: 8,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.RandomLen)
		})
	}
}

func TestSanitization_CallbackOnDestroy(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
	}{
		{
			name:   "Default to http://localhost/",
			sample: Model{},
			expect: "http://localhost/",
		},
		{
			name:   "Should has trailing slash",
			sample: Model{CallbackOnDestroy: "http://some-random.domain"},
			expect: "http://some-random.domain/",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.CallbackOnDestroy)
		})
	}
}

func TestSanitization_CallbackOnDestroyThisApp(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
	}{
		{
			name:   "Default to http://localhost",
			sample: Model{},
			expect: "http://localhost",
		},
		{
			name:   "Should not has trailing slash",
			sample: Model{CallbackOnDestroyThisApp: "http://some-random.domain/callback/url/"},
			expect: "http://some-random.domain/callback/url",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.CallbackOnDestroyThisApp)
		})
	}
}

func TestSanitizationLog(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
	}{
		{
			name:   "Log w /var/www/log/ should be /var/www/log/",
			sample: Model{LogDir: "/var/www/log/"},
			expect: "/var/www/log/",
		},
		{
			name:   "Log w /var/www/log should be /var/www/log/",
			sample: Model{LogDir: "/var/www/log"},
			expect: "/var/www/log/",
		},
		{
			name:   "Log w/o value should be ./log/",
			sample: Model{},
			expect: "./log/",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.sample.SanitizationLog()
			assert.Equal(t, tt.expect, tt.sample.LogDir)
		})
	}
}
