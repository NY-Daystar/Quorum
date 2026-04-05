package storage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emersion/go-message"
	"go.uber.org/zap"
)

// SaveMail save mail in output dir with generated id
func SaveMail(labelDir, filename string, raw []byte, log *zap.Logger) error {

	msg, err := message.Read(bytes.NewReader(raw))
	if err != nil {
		log.Sugar().Errorf("ReadBytes %v", err)
	}

	mailFolderPath := filepath.Join(labelDir, filename)
	os.MkdirAll(mailFolderPath, 0600)

	content := &MailContent{}
	ExtractParts(msg, content)

	// save inline if exists
	if len(content.InlineFiles) != 0 {
		inlineDir := filepath.Join(mailFolderPath, "inline")
		os.MkdirAll(inlineDir, 0600)

		for cid, data := range content.InlineFiles {
			path := filepath.Join(inlineDir, cid)
			os.WriteFile(path, data, 0600)
			content.HTML = strings.ReplaceAll(content.HTML, "cid:"+cid, "inline/"+cid)
		}
	}

	// save attachments if exists
	if len(content.Attachments) != 0 {
		attDir := filepath.Join(mailFolderPath, "attachments")
		os.MkdirAll(attDir, 0600)

		for name, data := range content.Attachments {
			os.WriteFile(filepath.Join(attDir, name), data, 0600)
		}
	}

	emlFile := filepath.Join(mailFolderPath, "mail.eml")
	err = SaveEML(emlFile, raw)
	if err != nil {
		return fmt.Errorf("SaveEML : %v", err)
	}

	htmlFile := filepath.Join(mailFolderPath, "mail.html")
	return SaveHTML(htmlFile, *msg, *content)
}

// SaveHTML convert eml file to html file
func SaveHTML(path string, message message.Entity, content MailContent) error {

	html := fmt.Sprintf(`
		<html><body>
		<h3>%s</h3>
		%s
		</body></html>`,
		message.Header.Get("Subject"),
		content.HTML,
	)

	return os.WriteFile(path, []byte(html), 0600)
}

// SaveEML save mail into eml format file
func SaveEML(path string, data []byte) error {
	return os.WriteFile(path, data, 0600)
}
