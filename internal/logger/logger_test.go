package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerWithDefaults(t *testing.T) {
	l := New(Options{})

	assert.NotNil(t, l)
	assert.NotNil(t, l.logger)
	assert.NotNil(t, l.styles)
	assert.False(t, l.noColor)
	assert.False(t, l.verbose)
}

func TestNewLoggerWithOptions(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{
		NoColor: true,
		Verbose: true,
		Output:  buf,
	})

	assert.True(t, l.noColor)
	assert.True(t, l.verbose)
}

func TestLoggerSuccess(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{Output: buf, NoColor: true})

	l.Success("operation completed")

	output := buf.String()
	assert.Contains(t, output, "[OK]")
	assert.Contains(t, output, "operation completed")
}

func TestLoggerError(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{Output: buf, NoColor: true})

	l.Error("something failed", "reason", "timeout")

	output := buf.String()
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "something failed")
	assert.Contains(t, output, "reason")
	assert.Contains(t, output, "timeout")
}

func TestLoggerWarning(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{Output: buf, NoColor: true})

	l.Warning("be careful")

	output := buf.String()
	assert.Contains(t, output, "[WARN]")
	assert.Contains(t, output, "be careful")
}

func TestLoggerInfo(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{Output: buf, NoColor: true})

	l.Info("some information")

	output := buf.String()
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "some information")
}

func TestLoggerDebugNotShownByDefault(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{Output: buf, NoColor: true})

	l.Debug("debug message")

	output := buf.String()
	assert.Empty(t, output)
}

func TestLoggerDebugShownWhenVerbose(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{Output: buf, NoColor: true, Verbose: true})

	l.Debug("debug message")

	output := buf.String()
	assert.Contains(t, output, "debug message")
}

func TestLoggerSetVerbose(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{Output: buf, NoColor: true})

	l.Debug("first debug")
	assert.Empty(t, buf.String())

	l.SetVerbose(true)
	l.Debug("second debug")
	assert.Contains(t, buf.String(), "second debug")
}

func TestLoggerSetNoColor(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{Output: buf, NoColor: false})

	l.Success("with color")
	outputWithColor := buf.String()
	buf.Reset()

	l.SetNoColor(true)
	l.Success("without color")
	outputWithoutColor := buf.String()

	assert.Contains(t, outputWithoutColor, "[OK]")
	assert.True(t, len(outputWithColor) != len(outputWithoutColor) ||
		!strings.Contains(outputWithColor, "[OK]"))
}

func TestDefaultLogger(t *testing.T) {
	assert.NotNil(t, Default())
}

func TestSetDefault(t *testing.T) {
	original := Default()
	defer SetDefault(original)

	buf := new(bytes.Buffer)
	newLogger := New(Options{Output: buf})
	SetDefault(newLogger)

	assert.Equal(t, newLogger, Default())
}

func TestPackageLevelFunctions(t *testing.T) {
	buf := new(bytes.Buffer)
	original := Default()
	defer SetDefault(original)

	SetDefault(New(Options{Output: buf, NoColor: true}))

	Success("pkg success")
	assert.Contains(t, buf.String(), "pkg success")
	buf.Reset()

	Error("pkg error")
	assert.Contains(t, buf.String(), "pkg error")
	buf.Reset()

	Warning("pkg warning")
	assert.Contains(t, buf.String(), "pkg warning")
	buf.Reset()

	Info("pkg info")
	assert.Contains(t, buf.String(), "pkg info")
	buf.Reset()

	SetVerbose(true)
	Debug("pkg debug")
	assert.Contains(t, buf.String(), "pkg debug")
}

func TestStylesCreation(t *testing.T) {
	stylesWithColor := createStyles(false)
	stylesNoColor := createStyles(true)

	assert.NotNil(t, stylesWithColor.Success)
	assert.NotNil(t, stylesWithColor.Error)
	assert.NotNil(t, stylesWithColor.Warning)
	assert.NotNil(t, stylesWithColor.Info)
	assert.NotNil(t, stylesWithColor.Debug)
	assert.NotNil(t, stylesWithColor.Bold)
	assert.NotNil(t, stylesWithColor.Muted)

	assert.NotNil(t, stylesNoColor.Success)
	assert.NotNil(t, stylesNoColor.Error)
}

func TestLoggerStyles(t *testing.T) {
	l := New(Options{})
	styles := l.Styles()

	assert.NotNil(t, styles)
	assert.NotNil(t, styles.Success)
	assert.NotNil(t, styles.Error)
}

func TestLoggerPrint(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(Options{Output: buf, NoColor: true})

	l.Print("plain message")

	output := buf.String()
	assert.Contains(t, output, "plain message")
}
