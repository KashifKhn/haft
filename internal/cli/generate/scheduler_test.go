package generate

import (
	"testing"

	"github.com/KashifKhn/haft/internal/detector"
	"github.com/stretchr/testify/assert"
)

func TestSchedulerCommand(t *testing.T) {
	cmd := newSchedulerCommand()

	assert.Equal(t, "scheduler <name>", cmd.Use)
	assert.Contains(t, cmd.Aliases, "sch")
	assert.Contains(t, cmd.Aliases, "scheduled")
	assert.Contains(t, cmd.Aliases, "task")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestSchedulerCommandFlags(t *testing.T) {
	cmd := newSchedulerCommand()

	pkgFlag := cmd.Flags().Lookup("package")
	assert.NotNil(t, pkgFlag)
	assert.Equal(t, "p", pkgFlag.Shorthand)

	cronFlag := cmd.Flags().Lookup("cron")
	assert.NotNil(t, cronFlag)

	rateFlag := cmd.Flags().Lookup("rate")
	assert.NotNil(t, rateFlag)

	delayFlag := cmd.Flags().Lookup("delay")
	assert.NotNil(t, delayFlag)

	initialFlag := cmd.Flags().Lookup("initial")
	assert.NotNil(t, initialFlag)

	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")
	assert.NotNil(t, noInteractiveFlag)

	refreshFlag := cmd.Flags().Lookup("refresh")
	assert.NotNil(t, refreshFlag)

	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)
}

func TestScheduleTypeOptionsDefinition(t *testing.T) {
	assert.Len(t, scheduleTypeOptions, 3)

	expectedTypes := []string{"cron", "fixedRate", "fixedDelay"}

	for i, expected := range expectedTypes {
		assert.Equal(t, expected, scheduleTypeOptions[i].Value)
		assert.NotEmpty(t, scheduleTypeOptions[i].Label)
		assert.NotEmpty(t, scheduleTypeOptions[i].Description)
	}
}

func TestGetSchedulerPackage(t *testing.T) {
	tests := []struct {
		name         string
		architecture detector.ArchitectureType
		basePackage  string
		expected     string
	}{
		{
			name:         "layered architecture",
			architecture: detector.ArchLayered,
			basePackage:  "com.example.app",
			expected:     "com.example.app.scheduler",
		},
		{
			name:         "feature architecture",
			architecture: detector.ArchFeature,
			basePackage:  "com.example.app",
			expected:     "com.example.app.common.scheduler",
		},
		{
			name:         "hexagonal architecture",
			architecture: detector.ArchHexagonal,
			basePackage:  "com.example.app",
			expected:     "com.example.app.infrastructure.scheduler",
		},
		{
			name:         "clean architecture",
			architecture: detector.ArchClean,
			basePackage:  "com.example.app",
			expected:     "com.example.app.infrastructure.scheduler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &detector.ProjectProfile{
				Architecture: tt.architecture,
				BasePackage:  tt.basePackage,
			}
			result := getSchedulerPackage(profile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildSchedulerTemplateData(t *testing.T) {
	profile := &detector.ProjectProfile{
		Architecture: detector.ArchLayered,
		BasePackage:  "com.example.app",
		Lombok:       detector.LombokProfile{Detected: true},
	}

	cfg := schedulerConfig{
		Name:           "Cleanup",
		ScheduleType:   "cron",
		CronExpression: "0 0 2 * * *",
	}

	data := buildSchedulerTemplateData(profile, cfg, "com.example.app.scheduler")

	assert.Equal(t, "Cleanup", data["Name"])
	assert.Equal(t, "com.example.app.scheduler", data["SchedulerPackage"])
	assert.Equal(t, "cron", data["ScheduleType"])
	assert.Equal(t, "0 0 2 * * *", data["CronExpression"])
	assert.Equal(t, true, data["HasLombok"])
	assert.Equal(t, "layered", data["Architecture"])
}

func TestBuildSchedulerTemplateData_FixedRate(t *testing.T) {
	profile := &detector.ProjectProfile{
		Architecture: detector.ArchLayered,
		BasePackage:  "com.example.app",
		Lombok:       detector.LombokProfile{Detected: false},
	}

	cfg := schedulerConfig{
		Name:         "Sync",
		ScheduleType: "fixedRate",
		FixedRate:    300000,
	}

	data := buildSchedulerTemplateData(profile, cfg, "com.example.app.scheduler")

	assert.Equal(t, "Sync", data["Name"])
	assert.Equal(t, "fixedRate", data["ScheduleType"])
	assert.Equal(t, int64(300000), data["FixedRate"])
	assert.Equal(t, false, data["HasLombok"])
}

func TestBuildSchedulerTemplateData_FixedDelay(t *testing.T) {
	profile := &detector.ProjectProfile{
		Architecture: detector.ArchFeature,
		BasePackage:  "com.example.app",
	}

	cfg := schedulerConfig{
		Name:         "Process",
		ScheduleType: "fixedDelay",
		FixedDelay:   60000,
		InitialDelay: 5000,
	}

	data := buildSchedulerTemplateData(profile, cfg, "com.example.app.common.scheduler")

	assert.Equal(t, "Process", data["Name"])
	assert.Equal(t, "fixedDelay", data["ScheduleType"])
	assert.Equal(t, int64(60000), data["FixedDelay"])
	assert.Equal(t, int64(5000), data["InitialDelay"])
	assert.Equal(t, "feature", data["Architecture"])
}

func TestSchedulerConfig_Defaults(t *testing.T) {
	cfg := schedulerConfig{}

	assert.Empty(t, cfg.Name)
	assert.Empty(t, cfg.BasePackage)
	assert.Empty(t, cfg.ScheduleType)
	assert.Empty(t, cfg.CronExpression)
	assert.Equal(t, int64(0), cfg.FixedRate)
	assert.Equal(t, int64(0), cfg.FixedDelay)
	assert.Equal(t, int64(0), cfg.InitialDelay)
}
