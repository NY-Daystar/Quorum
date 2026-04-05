package main

import (
	"context"
	"flag"
	"quorum/auth"
	"quorum/config"
	"quorum/mail"
	"quorum/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TODO faire un systeme d'options avec bubbletea (https://www.youtube.com/watch?v=16fWkt9OOPU)
// TODO rajouter des flags --since 2024-01-01
// TODO on est limité à 500 mails par labels voir comment faire sauter cette limite

func main() {
	log := initLogger()
	log.Sugar().Infof("Launch %s v%s", config.AppName, config.AppVersion)

	cfg, err := config.Load()
	if err != nil {
		log.Sugar().Warn("Can't load configuration, initialize...")
		log.Sugar().Debugf("error: %v", err)
		cfg = config.Init()
	}

	label := flag.String("label", "", "Gmail label to backup")
	anonymize := flag.Bool("anonymize", false, "Want to anonymize object mail")
	port := flag.Int64("port", config.AppPort, "Port to setup google service token, change redirect uril into crendential file")
	flag.Parse()
	cfg.Label = *label
	cfg.Anonymize = *anonymize
	cfg.Port = *port
	cfg.Save()

	log.Debug("configuration: ", zap.Any("configuration", cfg))

	client, err := auth.GetClient(cfg)
	if err != nil {
		log.Sugar().Fatalf("GetClient: %v", err)
	}

	ctx := context.Background()
	service, err := mail.NewService(ctx, client)
	if err != nil {
		log.Sugar().Fatalf("NewService %v", err)
	}

	var labels []mail.Label
	if cfg.Label == "" {
		labels, _ = mail.GetAllLabels(service)
	} else {
		labelID, err := mail.FindLabelID(service, cfg.Label)
		if err != nil {
			log.Sugar().Fatalf("FindLabelID: %v", err)
		}
		labels = append(labels, labelID)
	}

	mail.BackupMails(service, cfg, log, &labels)

	log.Sugar().Infof("✅ Backup done, mails are available into %v\n", cfg.OutputDir)
}

// initLogger create logger with zap librairy
func initLogger() *zap.Logger {
	logger := &lumberjack.Logger{
		Filename:   utils.GetLogsFile(), // File path
		MaxSize:    5,                   // 3 megabytes per files
		MaxBackups: 10,                  // 3 files before rotate
		MaxAge:     15,                  // 15 days
	}

	fileCore := utils.CreateFileLogger(logger)
	consoleCore := utils.CreateConsoleLogger()

	core := zapcore.NewTee(
		fileCore,
		consoleCore,
	)

	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer log.Sync() // flushes buffer, if any

	log.Debug("Zap logger set",
		zap.String("path", logger.Filename),
		zap.Int("filesize", logger.MaxSize), zap.Int("backupfile", logger.MaxBackups),
		zap.Int("fileage", logger.MaxAge),
	)

	return log
}
