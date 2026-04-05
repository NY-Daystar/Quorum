package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"quorum/config"
	"quorum/storage"
	"quorum/utils"

	"go.uber.org/zap"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Label represents label in gmail
type Label struct {
	ID   string
	Name string
}

// NewService init gmail service to get access
func NewService(ctx context.Context, client *http.Client) (*gmail.Service, error) {
	return gmail.NewService(ctx, option.WithHTTPClient(client))
}

// GetAllLabels return array with all gmail label
func GetAllLabels(s *gmail.Service) ([]Label, error) {
	labels, err := s.Users.Labels.List("me").Do()
	if err != nil {
		return nil, err
	}

	result := make([]Label, 0)

	for _, label := range labels.Labels {
		if label.Type == "system" {
			continue
		}
		result = append(result, Label{ID: label.Id, Name: label.Name})
	}

	return result, nil
}

// FindLabelID find specific id of label in gmail
func FindLabelID(s *gmail.Service, name string) (Label, error) {
	labels, err := s.Users.Labels.List("me").Do()
	if err != nil {
		return Label{}, err
	}

	for _, l := range labels.Labels {
		if l.Name == name {
			return Label{ID: l.Id, Name: l.Name}, nil
		}
	}

	return Label{}, fmt.Errorf("label %q not found", name)
}

// BackupMails get all messages in gmail account
func BackupMails(service *gmail.Service, cfg *config.Config, log *zap.Logger, labels *[]Label) {
	for idx, label := range *labels {
		msgs, err := getLabelMailList(service, label.ID, cfg)
		log.Sugar().Infof("%d - Label: %s (total mails:%d)", idx, label.Name, len(msgs))
		if err != nil {
			log.Sugar().Fatalf("ListMessage: %v", err)
		}

		for _, m := range msgs {
			raw, err := getMessage(service, m.Id, "raw")
			if err != nil {
				log.Sugar().Warnf("getMessage: %v", err)
				continue
			}
			rawString := raw.(string)

			data, err := base64.URLEncoding.DecodeString(rawString)
			if err != nil {
				log.Sugar().Warnf("DecodeString (%q) - %v", m.Id, err)
			}

			folderName := m.Id
			if !cfg.Anonymize {
				subject, date := getMessageMetadata(service, m.Id)
				parsedDate, err := utils.ParseDate(date)
				if err != nil {
					log.Sugar().Fatalf("BackupMails (ParseDate): %v", err)
				}
				safeSubject := utils.SanitizeFilename(subject)

				folderName = fmt.Sprintf("%s_%s", parsedDate, safeSubject)

				log.Sugar().Debugf("foldermail",
					zap.String("subject", subject),
					zap.String("date", date),
					zap.String("folderName", folderName),
				)
			}

			labelDir := filepath.Join(cfg.OutputDir, label.Name)
			err = storage.SaveMail(labelDir, folderName, data, log)
			if err != nil {
				log.Sugar().Fatalf("SaveMail (%q): %v", m.Id, err)
			}
		}
	}
}

// getLabelMailList get list of mail of specific label
func getLabelMailList(service *gmail.Service, labelID string, cfg *config.Config) ([]*gmail.Message, error) {
	user := "me"

	res, err := service.Users.Messages.List(user).LabelIds(labelID).MaxResults(cfg.MaxResults).Do() // TODO le maxResult est par defaut ilimité mettre un flag pour gérer

	if err != nil {
		return nil, err
	}
	return res.Messages, nil
}

// getMessage for each message get raw data ir full format
func getMessage(s *gmail.Service, id, format string) (interface{}, error) {
	msg, err := s.Users.Messages.Get("me", id).Format(format).Do()
	if err != nil {
		return "", err
	}
	if format == "raw" {
		return msg.Raw, nil
	}
	return msg, nil
}

// getMessage for each message get raw data ir full format
func getMessageMetadata(s *gmail.Service, id string) (string, string) {
	full, _ := getMessage(s, id, "full")
	msg := full.(*gmail.Message)
	return storage.GetHeader(msg.Payload.Headers, "Subject"), storage.GetHeader(msg.Payload.Headers, "Date")
}
