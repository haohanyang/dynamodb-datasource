package test

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func parseTime(datetime string, format string) (*time.Time, error) {
	cmd := exec.Command("node", "format.js", format)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(err.Error() + ":" + stderr.String())
	}

	layout := strings.Trim(stdout.String(), "\n")

	fmt.Println(layout)

	t, err := time.Parse(layout, datetime)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func TestFormatTime(t *testing.T) {
	t.Run("ISO 8601 - 1", func(t *testing.T) {
		ts, err := parseTime("2024-10-31T22:04:29+01:00", "YYYY-MM-DDTHH:mm:ssZ")
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, ts.Unix(), int64(1730408669))
	})

	t.Run("ISO 8601 - 2", func(t *testing.T) {
		ts, err := parseTime("2024-10-31T21:04:29.123Z", "YYYY-MM-DDTHH:mm:ss.SSS[Z]")
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, ts.UnixMilli(), int64(1730408669123))
	})

	t.Run("RFC1123", func(t *testing.T) {
		ts, err := parseTime("Thu, 31 Oct 2024 21:04:29 GMT", "ddd, DD MMM YYYY HH:mm:ss z")
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, ts.UnixMilli(), int64(1730408669000))
	})
}
