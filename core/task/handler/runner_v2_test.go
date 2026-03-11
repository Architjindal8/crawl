package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIpcDataRecords(t *testing.T) {
	line := `{"type":"data","payload":[{"test":"working"}],"ipc":true}`

	records, ok := parseIpcDataRecords(line)

	if !assert.True(t, ok) {
		return
	}
	if !assert.Len(t, records, 1) {
		return
	}
	assert.Equal(t, "working", records[0]["test"])
}

func TestParseIpcDataRecordsInvalid(t *testing.T) {
	records, ok := parseIpcDataRecords("plain log line")

	assert.False(t, ok)
	assert.Nil(t, records)
}
