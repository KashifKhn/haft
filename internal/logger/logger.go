package logger

import (
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	defaultLogger *Logger
)

type Logger struct {
	logger  *log.Logger
	noColor bool
	verbose bool
	output  io.Writer
	styles  *Styles
}

type Styles struct {
	Success lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Info    lipgloss.Style
	Debug   lipgloss.Style
	Bold    lipgloss.Style
	Muted   lipgloss.Style
}

type Options struct {
	NoColor bool
	Verbose bool
	Output  io.Writer
}

func init() {
	defaultLogger = New(Options{})
}

func New(opts Options) *Logger {
	output := opts.Output
	if output == nil {
		output = os.Stderr
	}

	l := log.NewWithOptions(output, log.Options{
		ReportTimestamp: false,
		ReportCaller:    false,
	})

	if opts.Verbose {
		l.SetLevel(log.DebugLevel)
	} else {
		l.SetLevel(log.InfoLevel)
	}

	styles := createStyles(opts.NoColor)
	customizeLogStyles(l, opts.NoColor)

	return &Logger{
		logger:  l,
		noColor: opts.NoColor,
		verbose: opts.Verbose,
		output:  output,
		styles:  styles,
	}
}

func createStyles(noColor bool) *Styles {
	if noColor {
		return &Styles{
			Success: lipgloss.NewStyle(),
			Error:   lipgloss.NewStyle(),
			Warning: lipgloss.NewStyle(),
			Info:    lipgloss.NewStyle(),
			Debug:   lipgloss.NewStyle(),
			Bold:    lipgloss.NewStyle(),
			Muted:   lipgloss.NewStyle(),
		}
	}

	return &Styles{
		Success: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		Error:   lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		Warning: lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
		Info:    lipgloss.NewStyle().Foreground(lipgloss.Color("12")),
		Debug:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Bold:    lipgloss.NewStyle().Bold(true),
		Muted:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
}

func customizeLogStyles(l *log.Logger, noColor bool) {
	if noColor {
		return
	}

	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Foreground(lipgloss.Color("8"))
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Foreground(lipgloss.Color("12"))
	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN").
		Foreground(lipgloss.Color("11"))
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Foreground(lipgloss.Color("9"))
	l.SetStyles(styles)
}

func (l *Logger) Success(msg string, keyvals ...interface{}) {
	prefix := "✓"
	if l.noColor {
		prefix = "[OK]"
	}
	l.logger.Info(l.styles.Success.Render(prefix)+" "+msg, keyvals...)
}

func (l *Logger) Error(msg string, keyvals ...interface{}) {
	prefix := "✗"
	if l.noColor {
		prefix = "[ERROR]"
	}
	l.logger.Error(l.styles.Error.Render(prefix)+" "+msg, keyvals...)
}

func (l *Logger) Warning(msg string, keyvals ...interface{}) {
	prefix := "⚠"
	if l.noColor {
		prefix = "[WARN]"
	}
	l.logger.Warn(l.styles.Warning.Render(prefix)+" "+msg, keyvals...)
}

func (l *Logger) Info(msg string, keyvals ...interface{}) {
	prefix := "ℹ"
	if l.noColor {
		prefix = "[INFO]"
	}
	l.logger.Info(l.styles.Info.Render(prefix)+" "+msg, keyvals...)
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.logger.Debug(msg, keyvals...)
}

func (l *Logger) Print(msg string) {
	l.logger.Print(msg)
}

func (l *Logger) Fatal(msg string, keyvals ...interface{}) {
	l.Error(msg, keyvals...)
	os.Exit(1)
}

func (l *Logger) SetVerbose(verbose bool) {
	l.verbose = verbose
	if verbose {
		l.logger.SetLevel(log.DebugLevel)
	} else {
		l.logger.SetLevel(log.InfoLevel)
	}
}

func (l *Logger) SetNoColor(noColor bool) {
	l.noColor = noColor
	l.styles = createStyles(noColor)
	customizeLogStyles(l.logger, noColor)
}

func (l *Logger) Styles() *Styles {
	return l.styles
}

func Default() *Logger {
	return defaultLogger
}

func SetDefault(l *Logger) {
	defaultLogger = l
}

func Success(msg string, keyvals ...interface{}) {
	defaultLogger.Success(msg, keyvals...)
}

func Error(msg string, keyvals ...interface{}) {
	defaultLogger.Error(msg, keyvals...)
}

func Warning(msg string, keyvals ...interface{}) {
	defaultLogger.Warning(msg, keyvals...)
}

func Info(msg string, keyvals ...interface{}) {
	defaultLogger.Info(msg, keyvals...)
}

func Debug(msg string, keyvals ...interface{}) {
	defaultLogger.Debug(msg, keyvals...)
}

func Print(msg string) {
	defaultLogger.Print(msg)
}

func SetVerbose(verbose bool) {
	defaultLogger.SetVerbose(verbose)
}

func SetNoColor(noColor bool) {
	defaultLogger.SetNoColor(noColor)
}
