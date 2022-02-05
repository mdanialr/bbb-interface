package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/kurvaid/bbb-interface/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeReader struct{}

func (f fakeReader) Read(_ []byte) (int, error) {
	return 0, errors.New("this should trigger error")
}

func TestSetup(t *testing.T) {
	fakeConfigFile :=
		`
env: prod
port: 6565
`
	var appConf config.Model

	t.Run("Using fake interface should return error", func(t *testing.T) {
		_, err := setup(&appConf, fakeReader{})
		require.Error(t, err)
	})

	// use manually created config file for this test only
	const tmpConfigPath = "/tmp/test-config-file.yml"

	f, err := os.Create(tmpConfigPath)
	require.NoError(t, err)
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalln("failed closing tmp config file:", err)
		}
	}()

	someRandomFile := bytes.Buffer{}
	someRandomFile.WriteString(fakeConfigFile)
	_, err = someRandomFile.WriteTo(f)
	require.NoError(t, err, "failed to write content to file")

	t.Cleanup(func() {
		// remove manually created config file
		if err := os.Remove(tmpConfigPath); err != nil {
			log.Fatalln("failed cleaning tmp config file:", err)
		}
		// remove leftover created log file
		if err := os.Remove("./log/app-log"); err != nil {
			log.Fatalln("failed cleaning log file:", err)
		}
	})

	// make sure ./log exist
	assert.NoError(t, os.MkdirAll("./log", 0777))

	t.Run("Success using real config file", func(t *testing.T) {
		f, err := os.ReadFile(tmpConfigPath)
		require.NoError(t, err)

		_, err = setup(&appConf, bytes.NewReader(f))
		require.NoError(t, err)
	})
	// end test that need manually created config file

	fakeConfigFile =
		`
env: prod
port: 6565
log: ./log
`
	t.Run("Success must exactly the same as in config file", func(t *testing.T) {
		_, err := setup(&appConf, bytes.NewBufferString(fakeConfigFile))
		require.NoError(t, err)

		assert.Equal(t, "localhost", appConf.Host)
		assert.Equal(t, strconv.Itoa(int(uint16(6565))), strconv.Itoa(int(appConf.PortNum)))
		assert.Equal(t, "./log/", appConf.LogDir)
	})

	fakeConfigFile =
		`
env: prod
port: 5050
log: /fake/dir
`

	t.Run("Log dir does not exist should return error", func(t *testing.T) {
		_, err := setup(&appConf, bytes.NewBufferString(fakeConfigFile))
		require.Error(t, err)
	})
}
