package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/KashifKhn/haft/internal/detector"
	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/output"
	"github.com/KashifKhn/haft/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type schedulerConfig struct {
	Name           string
	BasePackage    string
	ScheduleType   string
	CronExpression string
	FixedRate      int64
	FixedDelay     int64
	InitialDelay   int64
}

type scheduleTypeOption struct {
	Label       string
	Value       string
	Description string
}

var scheduleTypeOptions = []scheduleTypeOption{
	{"Cron Expression", "cron", "Run at specific times using cron syntax"},
	{"Fixed Rate", "fixedRate", "Run at fixed intervals (ms) from start"},
	{"Fixed Delay", "fixedDelay", "Run at fixed intervals (ms) after completion"},
}

func newSchedulerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduler <name>",
		Aliases: []string{"sch", "scheduled", "task"},
		Short:   "Generate a scheduled task class",
		Long: `Generate a Spring Boot @Scheduled task class.

The scheduler generator creates:
  - A @Component class with @Scheduled method
  - SchedulingConfig.java (if not exists) to enable scheduling

Schedule types:
  - cron: Run at specific times (e.g., "0 0 * * * *" for hourly)
  - fixedRate: Run every N milliseconds from method start
  - fixedDelay: Run N milliseconds after previous execution completes

Common cron expressions:
  - "0 0 * * * *"     Every hour
  - "0 0 8 * * *"     Every day at 8 AM
  - "0 0 0 * * MON"   Every Monday at midnight
  - "0 */15 * * * *"  Every 15 minutes`,
		Example: `  # Interactive mode
  haft generate scheduler cleanup
  haft g sch report

  # With cron expression
  haft generate scheduler cleanup --cron "0 0 2 * * *"

  # With fixed rate (every 5 minutes)
  haft generate scheduler sync --rate 300000

  # With fixed delay (1 minute after completion)
  haft generate scheduler process --delay 60000 --initial 5000

  # Non-interactive with all options
  haft generate scheduler report --cron "0 0 8 * * MON-FRI" --package com.example.app --no-interactive`,
		Args: cobra.MaximumNArgs(1),
		RunE: runScheduler,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from project)")
	cmd.Flags().String("cron", "", "Cron expression (e.g., '0 0 * * * *' for hourly)")
	cmd.Flags().Int64("rate", 0, "Fixed rate in milliseconds")
	cmd.Flags().Int64("delay", 0, "Fixed delay in milliseconds")
	cmd.Flags().Int64("initial", 0, "Initial delay in milliseconds (used with --delay)")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("refresh", false, "Force re-detection of project profile (ignore cache)")
	cmd.Flags().Bool("json", false, "Output result as JSON")

	return cmd
}

func runScheduler(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	forceRefresh, _ := cmd.Flags().GetBool("refresh")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	log := logger.Default()

	profile, err := DetectProjectProfileWithRefresh(forceRefresh)
	if err != nil {
		if noInteractive {
			if jsonOutput {
				return output.Error("DETECTION_ERROR", "Could not detect project profile", err.Error())
			}
			return fmt.Errorf("could not detect project profile: %w", err)
		}
		log.Warning("Could not detect project profile, using defaults")
		profile = &detector.ProjectProfile{
			Architecture: detector.ArchLayered,
		}
	}

	enrichProfileFromBuildFile(profile)

	if pkg, _ := cmd.Flags().GetString("package"); pkg != "" {
		profile.BasePackage = pkg
	}

	cfg := schedulerConfig{}

	if len(args) > 0 {
		cfg.Name = ToPascalCase(args[0])
	}

	cronExpr, _ := cmd.Flags().GetString("cron")
	fixedRate, _ := cmd.Flags().GetInt64("rate")
	fixedDelay, _ := cmd.Flags().GetInt64("delay")
	initialDelay, _ := cmd.Flags().GetInt64("initial")

	if cronExpr != "" {
		cfg.ScheduleType = "cron"
		cfg.CronExpression = cronExpr
	} else if fixedRate > 0 {
		cfg.ScheduleType = "fixedRate"
		cfg.FixedRate = fixedRate
	} else if fixedDelay > 0 {
		cfg.ScheduleType = "fixedDelay"
		cfg.FixedDelay = fixedDelay
		cfg.InitialDelay = initialDelay
	}

	if !noInteractive {
		wizardResult, err := runSchedulerWizard(profile.BasePackage, cfg)
		if err != nil {
			if jsonOutput {
				return output.Error("WIZARD_ERROR", "Wizard failed", err.Error())
			}
			return err
		}
		cfg = wizardResult
		if wizardResult.BasePackage != "" {
			profile.BasePackage = wizardResult.BasePackage
		}
	}

	if cfg.Name == "" {
		if jsonOutput {
			return output.Error("VALIDATION_ERROR", "Scheduler name is required")
		}
		return fmt.Errorf("scheduler name is required. Usage: haft generate scheduler <name>")
	}

	if profile.BasePackage == "" {
		if jsonOutput {
			return output.Error("VALIDATION_ERROR", "Base package could not be detected", "Use --package flag to specify it")
		}
		return fmt.Errorf("base package could not be detected. Use --package flag to specify it")
	}

	if cfg.ScheduleType == "" {
		cfg.ScheduleType = "cron"
		cfg.CronExpression = "0 0 * * * *"
	}

	return generateScheduler(profile, cfg, jsonOutput)
}

func runSchedulerWizard(currentPackage string, currentCfg schedulerConfig) (schedulerConfig, error) {
	cfg := currentCfg

	componentCfg := ComponentConfig{
		BasePackage: currentPackage,
		Name:        cfg.Name,
	}

	result, err := RunComponentWizard("Generate Scheduled Task", componentCfg, "Scheduler")
	if err != nil {
		return cfg, err
	}

	cfg.Name = result.Name
	cfg.BasePackage = result.BasePackage

	if cfg.ScheduleType == "" {
		schedType, err := runScheduleTypePicker()
		if err != nil {
			return cfg, err
		}
		cfg.ScheduleType = schedType

		switch cfg.ScheduleType {
		case "cron":
			cfg.CronExpression, err = runCronInput()
			if err != nil {
				return cfg, err
			}
		case "fixedRate":
			cfg.FixedRate, err = runRateInput("Enter fixed rate (milliseconds)", 60000)
			if err != nil {
				return cfg, err
			}
		case "fixedDelay":
			cfg.FixedDelay, err = runRateInput("Enter fixed delay (milliseconds)", 60000)
			if err != nil {
				return cfg, err
			}
			cfg.InitialDelay, err = runRateInput("Enter initial delay (milliseconds, 0 for none)", 0)
			if err != nil {
				return cfg, err
			}
		}
	}

	return cfg, nil
}

func runScheduleTypePicker() (string, error) {
	items := make([]components.SelectItem, len(scheduleTypeOptions))
	for i, opt := range scheduleTypeOptions {
		items[i] = components.SelectItem{
			Label:       opt.Label,
			Value:       opt.Value,
			Description: opt.Description,
		}
	}

	model := components.NewSelect(components.SelectConfig{
		Label: "Select schedule type",
		Items: items,
	})

	wrapper := scheduleTypePickerWrapper{model: model}
	p := tea.NewProgram(wrapper)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(scheduleTypePickerWrapper)
	if result.model.GoBack() {
		return "", fmt.Errorf("wizard cancelled")
	}

	return result.model.Value(), nil
}

type scheduleTypePickerWrapper struct {
	model components.SelectModel
}

func (w scheduleTypePickerWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w scheduleTypePickerWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return w, tea.Quit
		}
	}

	updated, cmd := w.model.Update(msg)
	w.model = updated

	if w.model.Submitted() || w.model.GoBack() {
		return w, tea.Quit
	}

	return w, cmd
}

func (w scheduleTypePickerWrapper) View() string {
	return w.model.View()
}

func runCronInput() (string, error) {
	model := components.NewTextInput(components.TextInputConfig{
		Label:       "Cron expression",
		Placeholder: "0 0 * * * *",
		Default:     "0 0 * * * *",
		HelpText:    "Format: sec min hour day month weekday (e.g., '0 0 8 * * MON-FRI' for 8 AM weekdays)",
	})

	wrapper := cronInputWrapper{model: model}
	p := tea.NewProgram(wrapper)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(cronInputWrapper)
	if result.model.GoBack() {
		return "", fmt.Errorf("wizard cancelled")
	}

	value := result.model.Value()
	if value == "" {
		value = "0 0 * * * *"
	}

	return value, nil
}

type cronInputWrapper struct {
	model components.TextInputModel
}

func (w cronInputWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w cronInputWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return w, tea.Quit
		}
	}

	updated, cmd := w.model.Update(msg)
	w.model = updated

	if w.model.Submitted() || w.model.GoBack() {
		return w, tea.Quit
	}

	return w, cmd
}

func (w cronInputWrapper) View() string {
	return w.model.View()
}

func runRateInput(label string, defaultValue int64) (int64, error) {
	model := components.NewTextInput(components.TextInputConfig{
		Label:       label,
		Placeholder: strconv.FormatInt(defaultValue, 10),
		Default:     strconv.FormatInt(defaultValue, 10),
		HelpText:    "Enter value in milliseconds",
	})

	wrapper := rateInputWrapper{model: model}
	p := tea.NewProgram(wrapper)
	finalModel, err := p.Run()
	if err != nil {
		return 0, err
	}

	result := finalModel.(rateInputWrapper)
	if result.model.GoBack() {
		return 0, fmt.Errorf("wizard cancelled")
	}

	value := result.model.Value()
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue, nil
	}

	return parsed, nil
}

type rateInputWrapper struct {
	model components.TextInputModel
}

func (w rateInputWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w rateInputWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return w, tea.Quit
		}
	}

	updated, cmd := w.model.Update(msg)
	w.model = updated

	if w.model.Submitted() || w.model.GoBack() {
		return w, tea.Quit
	}

	return w, cmd
}

func (w rateInputWrapper) View() string {
	return w.model.View()
}

func generateScheduler(profile *detector.ProjectProfile, cfg schedulerConfig, jsonOutput bool) error {
	log := logger.Default()
	fs := afero.NewOsFs()

	cwd, err := os.Getwd()
	if err != nil {
		if jsonOutput {
			return output.Error("DIRECTORY_ERROR", "Could not get current directory", err.Error())
		}
		return err
	}

	engine := generator.NewEngineWithLoader(fs, cwd)

	srcPath := FindSourcePath(cwd)
	if srcPath == "" {
		if jsonOutput {
			return output.Error("SOURCE_ERROR", "Could not find src/main/java directory")
		}
		return fmt.Errorf("could not find src/main/java directory")
	}

	schedulerPackage := getSchedulerPackage(profile)
	configPackage := getConfigPackage(profile)

	tracker := NewGenerateTracker("scheduler", "Scheduled Task")

	if !jsonOutput {
		log.Info("Generating scheduled task", "name", cfg.Name, "package", schedulerPackage)
	}

	configPackagePath := strings.ReplaceAll(configPackage, ".", string(os.PathSeparator))
	configBasePath := filepath.Join(srcPath, configPackagePath)
	configOutputPath := filepath.Join(configBasePath, "SchedulingConfig.java")
	configRelPath := FormatRelativePath(cwd, configOutputPath)

	if !engine.FileExists(configOutputPath) {
		configData := map[string]any{
			"ConfigPackage": configPackage,
		}

		if err := engine.RenderAndWrite("scheduler/SchedulingConfig.java.tmpl", configOutputPath, configData); err != nil {
			if jsonOutput {
				tracker.AddError(fmt.Sprintf("failed to generate SchedulingConfig: %s", err.Error()))
			} else {
				log.Warning("Failed to generate SchedulingConfig", "error", err)
			}
		} else {
			if !jsonOutput {
				log.Info("Created", "file", configRelPath)
			}
			tracker.AddGenerated(configRelPath)
		}
	}

	schedulerPackagePath := strings.ReplaceAll(schedulerPackage, ".", string(os.PathSeparator))
	schedulerBasePath := filepath.Join(srcPath, schedulerPackagePath)
	taskFileName := cfg.Name + "Task.java"
	taskOutputPath := filepath.Join(schedulerBasePath, taskFileName)
	taskRelPath := FormatRelativePath(cwd, taskOutputPath)

	if engine.FileExists(taskOutputPath) {
		if !jsonOutput {
			log.Warning("File exists, skipping", "file", taskRelPath)
		}
		tracker.AddSkipped(taskRelPath)
	} else {
		taskData := buildSchedulerTemplateData(profile, cfg, schedulerPackage)

		if err := engine.RenderAndWrite("scheduler/ScheduledTask.java.tmpl", taskOutputPath, taskData); err != nil {
			if jsonOutput {
				tracker.AddError(fmt.Sprintf("failed to generate %s: %s", taskFileName, err.Error()))
			} else {
				return fmt.Errorf("failed to generate %s: %w", taskFileName, err)
			}
		} else {
			if !jsonOutput {
				log.Info("Created", "file", taskRelPath)
			}
			tracker.AddGenerated(taskRelPath)
		}
	}

	if !jsonOutput && len(tracker.Generated) > 0 {
		log.Success(fmt.Sprintf("Generated %d scheduler files", len(tracker.Generated)))
		printSchedulerHelp(cfg)
	}

	return OutputGenerateResult(jsonOutput, tracker)
}

func getSchedulerPackage(profile *detector.ProjectProfile) string {
	switch profile.Architecture {
	case detector.ArchFeature:
		return profile.BasePackage + ".common.scheduler"
	case detector.ArchHexagonal:
		return profile.BasePackage + ".infrastructure.scheduler"
	case detector.ArchClean:
		return profile.BasePackage + ".infrastructure.scheduler"
	default:
		return profile.BasePackage + ".scheduler"
	}
}

func buildSchedulerTemplateData(profile *detector.ProjectProfile, cfg schedulerConfig, schedulerPackage string) map[string]any {
	return map[string]any{
		"Name":             cfg.Name,
		"SchedulerPackage": schedulerPackage,
		"ScheduleType":     cfg.ScheduleType,
		"CronExpression":   cfg.CronExpression,
		"FixedRate":        cfg.FixedRate,
		"FixedDelay":       cfg.FixedDelay,
		"InitialDelay":     cfg.InitialDelay,
		"HasLombok":        profile.Lombok.Detected,
		"Architecture":     string(profile.Architecture),
	}
}

func printSchedulerHelp(cfg schedulerConfig) {
	log := logger.Default()
	log.Info("")
	log.Info("Schedule configuration:")

	switch cfg.ScheduleType {
	case "cron":
		log.Info(fmt.Sprintf("  Cron: %s", cfg.CronExpression))
	case "fixedRate":
		log.Info(fmt.Sprintf("  Fixed Rate: %dms", cfg.FixedRate))
	case "fixedDelay":
		log.Info(fmt.Sprintf("  Fixed Delay: %dms (initial: %dms)", cfg.FixedDelay, cfg.InitialDelay))
	}

	log.Info("")
	log.Info("Common cron patterns:")
	log.Info("  0 0 * * * *      Every hour")
	log.Info("  0 0 8 * * *      Daily at 8 AM")
	log.Info("  0 0 0 * * MON    Every Monday at midnight")
	log.Info("  0 */15 * * * *   Every 15 minutes")
}
