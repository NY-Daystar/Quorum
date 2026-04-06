package console

import (
	"context"
	"fmt"
	"os"
	"quorum/auth"
	"quorum/config"
	"quorum/logger"
	"quorum/mail"
	"quorum/utils"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type screen int

const (
	menuScreen screen = iota
	backupSetupLabelScreen
	backupSetupAnonymizeScreen
	backupProgressScreen
	resultScreen
)

type action int

const (
	backupMailAction        action = 1
	sizeBackupAction        action = 2
	authenticateOAuthAction action = 3
	showCredentialsAction   action = 4
	showConfigAction        action = 5
	openLogsAction          action = 6
	quitAction              action = 7
)

// Choice type of choice
type Choice struct {
	Index action
	Name  string
}

// Model (tea.Model for bubbletea) represents all information of console
type Model struct {
	// For bubbletea
	Screen   screen
	Choices  []Choice
	Cursor   action
	Progress int
	Total    int
	Message  string
	Label    string

	// For Quorum
	Cfg *config.Config
	Log *logger.Log
}

// Start launch console and bubbletea program
func Start() Model {
	log := logger.Init()

	cfg, err := config.Load()
	if err != nil {
		log.Sugar().Warn("Can't load configuration, initialize...")
		log.Sugar().Debugf("error: %v", err)
		cfg = config.Init()
	}

	return Model{
		Screen: menuScreen,
		Cursor: 1,
		Choices: []Choice{
			{Index: backupMailAction, Name: "Backup mail"},
			{Index: sizeBackupAction, Name: "Count backup size"},
			{Index: authenticateOAuthAction, Name: "OAuth Authentication"},
			{Index: showCredentialsAction, Name: "Show credentials"},
			{Index: showConfigAction, Name: "Show configuration"},
			{Index: openLogsAction, Name: "Open logs"},
			{Index: quitAction, Name: "Quit"},
		},
		Log: logger.Init(),
		Cfg: cfg,
	}
}

// Init console interface
func (m Model) Init() tea.Cmd {
	m.Log.Debug("configuration: ", zap.Any("configuration", m.Cfg))
	return nil
}

// Update console interface after action
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.Screen {
		case menuScreen:
			keyStr := msg.String()
			m.Log.Debug("key", zap.String("key", keyStr))
			key, err := strconv.Atoi(keyStr)
			if err == nil {
				return m.runAction(action(key))
			}

			switch msg.String() {

			case "enter":
				return m.runAction(m.Cursor)
			case "q":
			case "ctrl+c":
				return m, tea.Quit

			case "up":
				if m.Cursor > 0 {
					m.Cursor--
				}

			case "down":
				if int(m.Cursor) < len(m.Choices)-1 {
					m.Cursor++
				}
			}

		case backupSetupLabelScreen:
			switch msg.Type {
			case tea.KeyCtrlC:
				m.Screen = menuScreen
				return m, nil
			case tea.KeyBackspace:
				if len(m.Cfg.Label) > 0 {
					m.Cfg.Label = m.Cfg.Label[:len(m.Cfg.Label)-1]
				}
			case tea.KeyEnter:
				m.Screen = backupSetupAnonymizeScreen
				m.Cfg.Save()
				return m, nil
			default:
				m.Cfg.Label += msg.String() // Get characters
			}

		case backupSetupAnonymizeScreen:
			switch msg.String() {
			case "enter":
				m.Screen = backupProgressScreen
				m.Cfg.Anonymize = m.Cursor == 0
				return m, m.runBackup()
			case "q":
				m.Screen = menuScreen
			case "right":
				m.Cursor = 1
			case "left":
				m.Cursor = 0
			}

		case backupProgressScreen:
			fmt.Println("SLT")
			if msg.String() == "enter" || msg.String() == "q" {
				m.Screen = menuScreen
			}

		case resultScreen:
			if msg.String() == "enter" || msg.String() == "q" {
				m.Screen = menuScreen
			}
		}

	case progressMsg:
		fmt.Println("progressMsg")
		m.Progress = msg.Current

		if m.Progress >= m.Total {
			m.Screen = resultScreen
			m.Message = "Backup done ✅"
		}
		return m, nil
	}

	return m, nil
}

// View show console interface
func (m Model) View() string {
	render := ""
	switch m.Screen {
	case menuScreen:
		render = fmt.Sprintf("\t\t\t%s (v%s)\n\t\t\tGmail Backup Tool\n\n",
			config.AppName, config.AppVersion) // TODO a changer avec de la couleur voir l'autre bibliotheque

		for _, c := range m.Choices {
			cursor := " "
			if m.Cursor == c.Index {
				cursor = ">"
			}
			render += fmt.Sprintf("%s %d - %s\n", cursor, c.Index, c.Name)
		}
		return render

	case backupSetupLabelScreen:
		return fmt.Sprintf(
			"Do you want to backup specific label (Default: \"ALL labels\"):\n\n> %s",
			m.Cfg.Label,
		)

	case backupSetupAnonymizeScreen:
		render += "Do you want to anonymize mail\n"

		opts := []string{"Yes", "No"}

		for idx, opt := range opts {
			cursor := " "
			if int(m.Cursor) == idx {
				cursor = ">"
			}
			render += fmt.Sprintf(" %s %s ", cursor, opt)
		}
		return render

	case backupProgressScreen:
		return fmt.Sprintf("Label: %v - Anonymize: %v\nBackup in progress...\n\n[%d/%d]\n",
			m.Cfg.Label, m.Cfg.Anonymize, m.Progress, m.Total)

	case resultScreen:
		return fmt.Sprintf("\n%s\n\n(press Enter or q to return)", m.Message)
	}

	return ""
}

// TODO a simplifier
func (m Model) runAction(choice action) (tea.Model, tea.Cmd) {
	switch choice {
	case backupMailAction:
		m.Screen = backupSetupLabelScreen
		return m, nil

	case sizeBackupAction:
		m.Message = fmt.Sprintf("Calculating size...\nSize is : %s", utils.CalculateFolderSize(m.Cfg.OutputDir))
		// TODO calculer le nombre de mail et de pièces jointes
		m.Screen = resultScreen
		return m, nil
	case authenticateOAuthAction:
		// TODO setup le port dans l'authentification
		// 	port := flag.Int64("port", config.AppPort, "Port to setup google service token, change redirect uril into crendential file")
		m.Log.Sugar().Infoln("Authenticate OAuth...") // TODO a retirer et mettre en message
		_, err := auth.GetClient(m.Cfg)
		if err != nil {
			m.Log.Sugar().Fatalf("GetClient: %v", err)
		}
		m.Log.Sugar().Infof("✅ Authenticated") // TODO a retirer et mettre en message
	case showCredentialsAction:
		data, err := os.ReadFile(m.Cfg.CredentialsFile)
		if err != nil {
			m.Message = "No credential file existing"
			return m, nil
		}
		m.Message = fmt.Sprintf("Showing credentials...\n%v", string(data))
		m.Screen = resultScreen
	case openLogsAction:
		data, _ := os.ReadFile(utils.GetLogsFile())
		m.Message = fmt.Sprintf("Showing logs...\n%v", string(data))
		m.Screen = resultScreen
	case showConfigAction:
		m.Message = fmt.Sprintf("Showing config...\n%v", m.Cfg.String())
		m.Screen = resultScreen
	case quitAction:
		m.Message = "Quitting application"
		return m, tea.Quit
	}

	m.Log.Sugar().Infoln(m.Message)
	return m, nil
}

type progressMsg struct {
	Current int
}

// TODO a eclater en envoyant plusieurs commande pour faire mail par mail
func (m Model) runBackup() tea.Cmd {
	m.Log.Sugar().Infoln("Running backup...")
	return func() tea.Msg {
		client, err := auth.GetClient(m.Cfg)
		if err != nil {
			m.Log.Sugar().Fatalf("GetClient: %v", err)
		}

		ctx := context.Background()
		service, err := mail.NewService(ctx, client)
		if err != nil {
			m.Log.Sugar().Fatalf("NewService %v", err)
		}

		var labels []mail.Label
		if m.Cfg.Label == "" {
			labels, _ = mail.GetAllLabels(service)
		} else {
			labelID, err := mail.FindLabelID(service, m.Cfg.Label)
			if err != nil {
				m.Log.Sugar().Fatalf("FindLabelID: %v", err)
			}
			labels = append(labels, labelID)
		}

		mail.BackupMails(service, m.Cfg, m.Log.ZapLogger, &labels)

		m.Log.Sugar().Infof("✅ Backup done, mails are available into %v\n", m.Cfg.OutputDir)
		return progressMsg{Current: 1} // TODO a travailler label par label ou mail par mail
	}
}
