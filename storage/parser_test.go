package storage

import (
	"bytes"
	"testing"

	"github.com/emersion/go-message"
)

func TestExtractHTML(t *testing.T) {
	raw := "Content-Type: text/html\r\n\r\n<h1>Hello</h1>"

	m, _ := message.Read(bytes.NewReader([]byte(raw)))

	content := &MailContent{}
	ExtractParts(m, content)

	if content.HTML != "<h1>Hello</h1>" {
		t.Errorf("HTML incorrect")
	}
}
