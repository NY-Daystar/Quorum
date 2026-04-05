package storage

import (
	"io"
	"strings"

	"github.com/emersion/go-message"
	"google.golang.org/api/gmail/v1"
)

// MailContent represents mail structure
type MailContent struct {
	HTML        string
	InlineFiles map[string][]byte
	Attachments map[string][]byte
}

// ExtractParts to get each part of a mail
func ExtractParts(e *message.Entity, content *MailContent) {
	if isMultipart(e) {
		processMultipart(e, content)
		return
	}

	processSinglePart(e, content)
}

func isMultipart(e *message.Entity) bool {
	return e.MultipartReader() != nil
}

func processMultipart(e *message.Entity, content *MailContent) {
	mr := e.MultipartReader()

	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}
		ExtractParts(part, content)
	}
}

func processSinglePart(e *message.Entity, content *MailContent) {
	ct, params, _ := e.Header.ContentType()
	data, _ := io.ReadAll(e.Body)

	if isHTML(ct) {
		content.HTML = string(data)
		return
	}

	if cid := getContentID(e); cid != "" {
		addInline(content, cid, data)
		return
	}

	if filename := getFilename(params); filename != "" {
		addAttachment(content, filename, data)
	}
}

func isHTML(ct string) bool {
	return ct == "text/html"
}

func getContentID(e *message.Entity) string {
	cid := e.Header.Get("Content-Id")
	return strings.Trim(cid, "<>")
}

func getFilename(params map[string]string) string {
	if name, ok := params["name"]; ok {
		return name
	}
	return ""
}

func addInline(content *MailContent, cid string, data []byte) {
	if content.InlineFiles == nil {
		content.InlineFiles = map[string][]byte{}
	}
	content.InlineFiles[cid] = data
}

func addAttachment(content *MailContent, name string, data []byte) {
	if content.Attachments == nil {
		content.Attachments = map[string][]byte{}
	}
	content.Attachments[name] = data
}

func GetHeader(headers []*gmail.MessagePartHeader, name string) string {
	for _, h := range headers {
		if h.Name == name {
			return h.Value
		}
	}
	return ""
}
