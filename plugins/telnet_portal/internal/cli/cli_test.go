package cli

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	b := bytes.NewBufferString("")
	CLI.SetOut(b)

	CLI.SetArgs([]string{"--help"})
	err := CLI.Execute()

	if err != nil {
		t.Fatalf("error executing CLI: %s", err)
	}

	out, err := ioutil.ReadAll(b)

	if err != nil {
		t.Fatalf("error reading output: %s", err)
	}

	assert.Equal(t, "start the Mjolnir Telnet Portal\n\n", string(out))

}
