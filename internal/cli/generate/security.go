package generate

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/detector"
	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/output"
	"github.com/KashifKhn/haft/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type SecurityType string

const (
	SecurityJWT     SecurityType = "jwt"
	SecuritySession SecurityType = "session"
	SecurityOAuth2  SecurityType = "oauth2"
)

type securityConfig struct {
	BasePackage      string
	SecurityTypes    []SecurityType
	GenerateEntities bool
	UserEntityName   string
}

type securityDependency struct {
	Name       string
	GroupId    string
	ArtifactId string
	Version    string
	Required   []SecurityType
}

var securityDependencies = []securityDependency{
	{"Spring Security", "org.springframework.boot", "spring-boot-starter-security", "", nil},
	{"Spring Data JPA", "org.springframework.boot", "spring-boot-starter-data-jpa", "", nil},
	{"JJWT API", "io.jsonwebtoken", "jjwt-api", "0.12.6", []SecurityType{SecurityJWT}},
	{"JJWT Impl", "io.jsonwebtoken", "jjwt-impl", "0.12.6", []SecurityType{SecurityJWT}},
	{"JJWT Jackson", "io.jsonwebtoken", "jjwt-jackson", "0.12.6", []SecurityType{SecurityJWT}},
	{"OAuth2 Client", "org.springframework.boot", "spring-boot-starter-oauth2-client", "", []SecurityType{SecurityOAuth2}},
}

var userEntityNames = []string{"User", "AppUser", "Account", "Member", "Principal", "UserEntity", "ApplicationUser"}
var roleEntityNames = []string{"Role", "Authority", "Permission", "UserRole"}

func newSecurityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security",
		Aliases: []string{"sec", "auth"},
		Short:   "Generate security configuration",
		Long: `Generate Spring Security configuration with authentication support.

Available authentication types:
  - JWT: Stateless token-based authentication (recommended for APIs)
  - Session: Traditional session-based authentication (for web apps)
  - OAuth2: OAuth2/OpenID Connect authentication (Google, GitHub, etc.)

The generator will:
  - Check for required dependencies and prompt to add missing ones
  - Detect existing User/Role entities or generate new ones
  - Generate appropriate configuration based on your project structure

Files generated for JWT:
  - SecurityConfig.java — Spring Security configuration
  - JwtUtil.java — JWT token utilities
  - JwtAuthenticationFilter.java — Token validation filter
  - AuthenticationController.java — Login/Register endpoints
  - AuthRequest.java, AuthResponse.java, RegisterRequest.java — DTOs
  - CustomUserDetailsService.java — User loading service
  - User.java, Role.java — Entities (if not detected)
  - UserRepository.java, RoleRepository.java — Repositories (if not detected)`,
		Example: `  # Interactive picker to select authentication type
  haft generate security
  haft g sec

  # Generate JWT authentication
  haft generate security --jwt

  # Generate session-based authentication
  haft generate security --session

  # Generate OAuth2 authentication
  haft generate security --oauth2

  # Generate all authentication types
  haft generate security --all

  # Override base package
  haft generate security --jwt --package com.example.app`,
		RunE: runSecurity,
	}

	cmd.Flags().StringP("package", "p", "", "Base package (auto-detected from project)")
	cmd.Flags().Bool("jwt", false, "Generate JWT authentication")
	cmd.Flags().Bool("session", false, "Generate session-based authentication")
	cmd.Flags().Bool("oauth2", false, "Generate OAuth2 authentication")
	cmd.Flags().Bool("all", false, "Generate all authentication types")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive wizard")
	cmd.Flags().Bool("skip-entities", false, "Skip User/Role entity generation even if not detected")
	cmd.Flags().Bool("refresh", false, "Force re-detection of project profile (ignore cache)")
	cmd.Flags().Bool("json", false, "Output result as JSON")

	return cmd
}

func runSecurity(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	includeAll, _ := cmd.Flags().GetBool("all")
	jwtFlag, _ := cmd.Flags().GetBool("jwt")
	sessionFlag, _ := cmd.Flags().GetBool("session")
	oauth2Flag, _ := cmd.Flags().GetBool("oauth2")
	skipEntities, _ := cmd.Flags().GetBool("skip-entities")
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

	cfg := securityConfig{}

	if includeAll {
		cfg.SecurityTypes = []SecurityType{SecurityJWT, SecuritySession, SecurityOAuth2}
	} else if jwtFlag || sessionFlag || oauth2Flag {
		if jwtFlag {
			cfg.SecurityTypes = append(cfg.SecurityTypes, SecurityJWT)
		}
		if sessionFlag {
			cfg.SecurityTypes = append(cfg.SecurityTypes, SecuritySession)
		}
		if oauth2Flag {
			cfg.SecurityTypes = append(cfg.SecurityTypes, SecurityOAuth2)
		}
	}

	if !noInteractive {
		wizardCfg, err := runSecurityWizard(profile.BasePackage, len(cfg.SecurityTypes) > 0)
		if err != nil {
			if jsonOutput {
				return output.Error("WIZARD_ERROR", "Wizard failed", err.Error())
			}
			return err
		}
		if wizardCfg.BasePackage != "" {
			profile.BasePackage = wizardCfg.BasePackage
		}
		if len(cfg.SecurityTypes) == 0 {
			cfg.SecurityTypes = wizardCfg.SecurityTypes
		}
	}

	if profile.BasePackage == "" {
		if jsonOutput {
			return output.Error("VALIDATION_ERROR", "Base package could not be detected", "Use --package flag to specify it")
		}
		return fmt.Errorf("base package could not be detected. Use --package flag to specify it (e.g., --package com.example.myapp)")
	}

	if len(cfg.SecurityTypes) == 0 {
		if noInteractive {
			if jsonOutput {
				return output.Error("VALIDATION_ERROR", "No authentication type specified", "Use --jwt, --session, --oauth2, or --all")
			}
			return fmt.Errorf("specify authentication type with --jwt, --session, --oauth2, or --all")
		}
		if jsonOutput {
			return output.Success(output.GenerateOutput{
				Results:        []output.GenerateResult{},
				TotalGenerated: 0,
				TotalSkipped:   0,
			})
		}
		log.Info("No authentication type selected")
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		if jsonOutput {
			return output.Error("DIRECTORY_ERROR", "Could not get current directory", err.Error())
		}
		return err
	}

	fs := afero.NewOsFs()
	missingDeps, err := checkSecurityDependencies(cwd, fs, cfg.SecurityTypes)
	if err != nil {
		if jsonOutput {
			return output.Error("DEPENDENCY_ERROR", "Could not check dependencies", err.Error())
		}
		return err
	}

	if len(missingDeps) > 0 {
		if noInteractive {
			if jsonOutput {
				return output.Error("DEPENDENCY_ERROR", "Missing required dependencies", formatMissingDeps(missingDeps))
			}
			return fmt.Errorf("missing required dependencies: %v", formatMissingDeps(missingDeps))
		}

		shouldAdd, err := promptAddDependencies(missingDeps)
		if err != nil {
			if jsonOutput {
				return output.Error("WIZARD_ERROR", "Dependency prompt failed", err.Error())
			}
			return err
		}

		if shouldAdd {
			if err := addSecurityDependencies(cwd, fs, missingDeps); err != nil {
				if jsonOutput {
					return output.Error("DEPENDENCY_ERROR", "Failed to add dependencies", err.Error())
				}
				return err
			}
		} else {
			log.Warning("Skipping dependency installation. Generated code may not compile.")
		}
	}

	userEntity, _ := detectUserEntity(cwd, fs, profile.BasePackage)
	cfg.GenerateEntities = userEntity == "" && !skipEntities
	if userEntity != "" {
		cfg.UserEntityName = userEntity
	} else {
		cfg.UserEntityName = "User"
	}

	if cfg.GenerateEntities && !noInteractive {
		shouldGenerate, err := promptGenerateEntities()
		if err != nil {
			if jsonOutput {
				return output.Error("WIZARD_ERROR", "Entity prompt failed", err.Error())
			}
			return err
		}
		cfg.GenerateEntities = shouldGenerate
	}

	return generateSecurity(profile, cfg, jsonOutput)
}

func runSecurityWizard(currentPackage string, skipTypePicker bool) (securityConfig, error) {
	cfg := securityConfig{}

	componentCfg := ComponentConfig{
		BasePackage: currentPackage,
		Name:        "Security",
	}

	result, err := RunComponentWizard("Generate Security Configuration", componentCfg, "Security")
	if err != nil {
		return cfg, err
	}
	cfg.BasePackage = result.BasePackage

	if skipTypePicker {
		return cfg, nil
	}

	cfg.SecurityTypes, err = runSecurityTypePicker()
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func runSecurityTypePicker() ([]SecurityType, error) {
	items := []components.MultiSelectItem{
		{Label: "JWT — Stateless token-based authentication (recommended for APIs)", Value: "jwt", Selected: false},
		{Label: "Session — Traditional session-based authentication (for web apps)", Value: "session", Selected: false},
		{Label: "OAuth2 — OAuth2/OpenID Connect (Google, GitHub, etc.)", Value: "oauth2", Selected: false},
	}

	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label:    "Select authentication types to generate",
		Items:    items,
		Required: true,
	})

	wrapper := securityMultiSelectWrapper{model: model}
	p := tea.NewProgram(wrapper)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(securityMultiSelectWrapper)
	if result.model.GoBack() {
		return nil, fmt.Errorf("wizard cancelled")
	}

	values := result.model.Values()
	var types []SecurityType
	for _, v := range values {
		types = append(types, SecurityType(v))
	}

	return types, nil
}

type securityMultiSelectWrapper struct {
	model components.MultiSelectModel
}

func (w securityMultiSelectWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w securityMultiSelectWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (w securityMultiSelectWrapper) View() string {
	return w.model.View()
}

func checkSecurityDependencies(cwd string, fs afero.Fs, types []SecurityType) ([]securityDependency, error) {
	result, err := buildtool.Detect(cwd, fs)
	if err != nil {
		return nil, fmt.Errorf("could not find build file: %w", err)
	}

	project, err := result.Parser.Parse(result.FilePath)
	if err != nil {
		return nil, fmt.Errorf("could not parse build file: %w", err)
	}

	var missing []securityDependency
	for _, dep := range securityDependencies {
		if dep.Required != nil && !containsAnySecurityType(types, dep.Required) {
			continue
		}

		if !result.Parser.HasDependency(project, dep.GroupId, dep.ArtifactId) {
			missing = append(missing, dep)
		}
	}

	return missing, nil
}

func containsAnySecurityType(haystack []SecurityType, needles []SecurityType) bool {
	for _, needle := range needles {
		for _, h := range haystack {
			if h == needle {
				return true
			}
		}
	}
	return false
}

func formatMissingDeps(deps []securityDependency) string {
	var names []string
	for _, dep := range deps {
		names = append(names, dep.Name)
	}
	return strings.Join(names, ", ")
}

func promptAddDependencies(deps []securityDependency) (bool, error) {
	log := logger.Default()
	log.Warning("Missing required dependencies:")
	for _, dep := range deps {
		log.Info("  - " + dep.Name + " (" + dep.GroupId + ":" + dep.ArtifactId + ")")
	}

	items := []components.SelectItem{
		{Label: "Yes, add missing dependencies", Value: "yes"},
		{Label: "No, continue without adding", Value: "no"},
	}

	model := components.NewSelect(components.SelectConfig{
		Label: "Add missing dependencies to your project?",
		Items: items,
	})

	wrapper := selectWrapper{model: model}
	p := tea.NewProgram(wrapper)
	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	result := finalModel.(selectWrapper)
	if result.model.GoBack() {
		return false, fmt.Errorf("wizard cancelled")
	}

	return result.model.Value() == "yes", nil
}

type selectWrapper struct {
	model components.SelectModel
}

func (w selectWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w selectWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (w selectWrapper) View() string {
	return w.model.View()
}

func addSecurityDependencies(cwd string, fs afero.Fs, deps []securityDependency) error {
	log := logger.Default()

	result, err := buildtool.Detect(cwd, fs)
	if err != nil {
		return err
	}

	project, err := result.Parser.Parse(result.FilePath)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		result.Parser.AddDependency(project, buildtool.Dependency{
			GroupId:    dep.GroupId,
			ArtifactId: dep.ArtifactId,
			Version:    dep.Version,
		})
		log.Success("Added", "dependency", dep.Name)
	}

	if err := result.Parser.Write(result.FilePath, project); err != nil {
		return fmt.Errorf("could not write build file: %w", err)
	}

	return nil
}

func detectUserEntity(cwd string, fs afero.Fs, basePackage string) (string, string) {
	srcPath := FindSourcePath(cwd)
	if srcPath == "" {
		return "", ""
	}

	packagePath := strings.ReplaceAll(basePackage, ".", string(os.PathSeparator))
	searchPaths := []string{
		filepath.Join(srcPath, packagePath, "entity"),
		filepath.Join(srcPath, packagePath, "model"),
		filepath.Join(srcPath, packagePath, "domain"),
		filepath.Join(srcPath, packagePath, "user"),
		filepath.Join(srcPath, packagePath, "auth"),
		filepath.Join(srcPath, packagePath),
	}

	for _, searchPath := range searchPaths {
		for _, entityName := range userEntityNames {
			filePath := filepath.Join(searchPath, entityName+".java")
			if exists, _ := afero.Exists(fs, filePath); exists {
				return entityName, searchPath
			}
		}
	}

	return "", ""
}

func promptGenerateEntities() (bool, error) {
	log := logger.Default()
	log.Warning("No User entity detected in your project")

	items := []components.SelectItem{
		{Label: "Yes, generate User and Role entities", Value: "yes"},
		{Label: "No, I'll create them manually", Value: "no"},
	}

	model := components.NewSelect(components.SelectConfig{
		Label: "Generate User and Role entities?",
		Items: items,
	})

	wrapper := selectWrapper{model: model}
	p := tea.NewProgram(wrapper)
	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	result := finalModel.(selectWrapper)
	if result.model.GoBack() {
		return false, fmt.Errorf("wizard cancelled")
	}

	return result.model.Value() == "yes", nil
}

func generateSecurity(profile *detector.ProjectProfile, cfg securityConfig, jsonOutput bool) error {
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

	securityPackage := getSecurityPackage(profile)
	userEntityPackage := getUserEntityPackage(profile)
	userRepositoryPackage := getUserRepositoryPackage(profile)

	data := buildSecurityTemplateData(profile, securityPackage, userEntityPackage, userRepositoryPackage, cfg)

	tracker := NewGenerateTracker("security", "SecurityConfig")

	if !jsonOutput {
		log.Info("Generating security configuration", "package", securityPackage)
	}

	for _, secType := range cfg.SecurityTypes {
		count, skipped, err := generateSecurityTypeTracked(engine, srcPath, secType, data, cfg, cwd, tracker, jsonOutput)
		if err != nil {
			if jsonOutput {
				tracker.AddError(err.Error())
				continue
			}
			return err
		}
		_ = count
		_ = skipped
	}

	if cfg.GenerateEntities {
		_, _, err := generateSecurityEntitiesTracked(engine, srcPath, data, cwd, tracker, jsonOutput)
		if err != nil {
			if jsonOutput {
				tracker.AddError(err.Error())
			} else {
				return err
			}
		}
	}

	if !jsonOutput {
		if len(tracker.Generated) > 0 {
			log.Success(fmt.Sprintf("Generated %d security files", len(tracker.Generated)))
		}
		if len(tracker.Skipped) > 0 {
			log.Info(fmt.Sprintf("Skipped %d existing files", len(tracker.Skipped)))
		}
		printSecurityInstructions(cfg.SecurityTypes)
	}

	return OutputGenerateResult(jsonOutput, tracker)
}

func generateSecurityType(engine *generator.Engine, srcPath string, secType SecurityType, data map[string]any, cfg securityConfig, cwd string) (int, int, error) {
	log := logger.Default()

	securityPackage := data["SecurityPackage"].(string)
	packagePath := strings.ReplaceAll(securityPackage, ".", string(os.PathSeparator))
	basePath := filepath.Join(srcPath, packagePath)

	templates := getSecurityTemplates(secType)

	generatedCount := 0
	skippedCount := 0

	for _, t := range templates {
		if t.conditional != "" {
			if val, ok := data[t.conditional].(bool); !ok || !val {
				continue
			}
		}

		outputPath := filepath.Join(basePath, t.fileName)

		if engine.FileExists(outputPath) {
			log.Warning("File exists, skipping", "file", FormatRelativePath(cwd, outputPath))
			skippedCount++
			continue
		}

		if err := engine.RenderAndWrite(t.template, outputPath, data); err != nil {
			return generatedCount, skippedCount, fmt.Errorf("failed to generate %s: %w", t.fileName, err)
		}

		log.Info("Created", "file", FormatRelativePath(cwd, outputPath))
		generatedCount++
	}

	return generatedCount, skippedCount, nil
}

func generateSecurityTypeTracked(engine *generator.Engine, srcPath string, secType SecurityType, data map[string]any, cfg securityConfig, cwd string, tracker *GenerateTracker, jsonOutput bool) (int, int, error) {
	log := logger.Default()

	securityPackage := data["SecurityPackage"].(string)
	packagePath := strings.ReplaceAll(securityPackage, ".", string(os.PathSeparator))
	basePath := filepath.Join(srcPath, packagePath)

	templates := getSecurityTemplates(secType)

	generatedCount := 0
	skippedCount := 0

	for _, t := range templates {
		if t.conditional != "" {
			if val, ok := data[t.conditional].(bool); !ok || !val {
				continue
			}
		}

		outputPath := filepath.Join(basePath, t.fileName)
		relPath := FormatRelativePath(cwd, outputPath)

		if engine.FileExists(outputPath) {
			if !jsonOutput {
				log.Warning("File exists, skipping", "file", relPath)
			}
			tracker.AddSkipped(relPath)
			skippedCount++
			continue
		}

		if err := engine.RenderAndWrite(t.template, outputPath, data); err != nil {
			if jsonOutput {
				tracker.AddError(fmt.Sprintf("failed to generate %s: %s", t.fileName, err.Error()))
				continue
			}
			return generatedCount, skippedCount, fmt.Errorf("failed to generate %s: %w", t.fileName, err)
		}

		if !jsonOutput {
			log.Info("Created", "file", relPath)
		}
		tracker.AddGenerated(relPath)
		generatedCount++
	}

	return generatedCount, skippedCount, nil
}

type templateInfo struct {
	template    string
	fileName    string
	conditional string
}

func getSecurityTemplates(secType SecurityType) []templateInfo {
	switch secType {
	case SecurityJWT:
		return []templateInfo{
			{"security/jwt/SecurityConfig.java.tmpl", "SecurityConfig.java", ""},
			{"security/jwt/JwtUtil.java.tmpl", "JwtUtil.java", ""},
			{"security/jwt/JwtAuthenticationFilter.java.tmpl", "JwtAuthenticationFilter.java", ""},
			{"security/jwt/AuthenticationController.java.tmpl", "AuthenticationController.java", ""},
			{"security/jwt/AuthRequest.java.tmpl", "AuthRequest.java", ""},
			{"security/jwt/AuthResponse.java.tmpl", "AuthResponse.java", ""},
			{"security/jwt/RegisterRequest.java.tmpl", "RegisterRequest.java", ""},
			{"security/jwt/RefreshTokenRequest.java.tmpl", "RefreshTokenRequest.java", ""},
			{"security/jwt/CustomUserDetailsService.java.tmpl", "CustomUserDetailsService.java", ""},
		}
	case SecuritySession:
		return []templateInfo{
			{"security/session/SecurityConfig.java.tmpl", "SecurityConfig.java", ""},
			{"security/session/CustomUserDetailsService.java.tmpl", "CustomUserDetailsService.java", ""},
			{"security/session/AuthController.java.tmpl", "AuthController.java", ""},
			{"security/session/RegisterRequest.java.tmpl", "RegisterRequest.java", ""},
		}
	case SecurityOAuth2:
		return []templateInfo{
			{"security/oauth2/SecurityConfig.java.tmpl", "SecurityConfig.java", ""},
			{"security/oauth2/OAuth2UserService.java.tmpl", "OAuth2UserService.java", ""},
			{"security/oauth2/OAuth2SuccessHandler.java.tmpl", "OAuth2SuccessHandler.java", ""},
			{"security/oauth2/OAuth2UserPrincipal.java.tmpl", "OAuth2UserPrincipal.java", ""},
		}
	default:
		return nil
	}
}

func generateSecurityEntities(engine *generator.Engine, srcPath string, data map[string]any, cwd string) (int, int, error) {
	log := logger.Default()

	userEntityPackage := data["UserEntityPackage"].(string)
	userRepoPackage := data["UserRepositoryPackage"].(string)

	entityPath := strings.ReplaceAll(userEntityPackage, ".", string(os.PathSeparator))
	repoPath := strings.ReplaceAll(userRepoPackage, ".", string(os.PathSeparator))

	entityBasePath := filepath.Join(srcPath, entityPath)
	repoBasePath := filepath.Join(srcPath, repoPath)

	templates := []struct {
		template   string
		fileName   string
		outputPath string
	}{
		{"security/jwt/User.java.tmpl", data["UserEntityName"].(string) + ".java", entityBasePath},
		{"security/jwt/Role.java.tmpl", "Role.java", entityBasePath},
		{"security/jwt/UserRepository.java.tmpl", data["UserEntityName"].(string) + "Repository.java", repoBasePath},
		{"security/jwt/RoleRepository.java.tmpl", "RoleRepository.java", repoBasePath},
	}

	generatedCount := 0
	skippedCount := 0

	for _, t := range templates {
		outputPath := filepath.Join(t.outputPath, t.fileName)

		if engine.FileExists(outputPath) {
			log.Warning("File exists, skipping", "file", FormatRelativePath(cwd, outputPath))
			skippedCount++
			continue
		}

		if err := engine.RenderAndWrite(t.template, outputPath, data); err != nil {
			return generatedCount, skippedCount, fmt.Errorf("failed to generate %s: %w", t.fileName, err)
		}

		log.Info("Created", "file", FormatRelativePath(cwd, outputPath))
		generatedCount++
	}

	return generatedCount, skippedCount, nil
}

func generateSecurityEntitiesTracked(engine *generator.Engine, srcPath string, data map[string]any, cwd string, tracker *GenerateTracker, jsonOutput bool) (int, int, error) {
	log := logger.Default()

	userEntityPackage := data["UserEntityPackage"].(string)
	userRepoPackage := data["UserRepositoryPackage"].(string)

	entityPath := strings.ReplaceAll(userEntityPackage, ".", string(os.PathSeparator))
	repoPath := strings.ReplaceAll(userRepoPackage, ".", string(os.PathSeparator))

	entityBasePath := filepath.Join(srcPath, entityPath)
	repoBasePath := filepath.Join(srcPath, repoPath)

	templates := []struct {
		template   string
		fileName   string
		outputPath string
	}{
		{"security/jwt/User.java.tmpl", data["UserEntityName"].(string) + ".java", entityBasePath},
		{"security/jwt/Role.java.tmpl", "Role.java", entityBasePath},
		{"security/jwt/UserRepository.java.tmpl", data["UserEntityName"].(string) + "Repository.java", repoBasePath},
		{"security/jwt/RoleRepository.java.tmpl", "RoleRepository.java", repoBasePath},
	}

	generatedCount := 0
	skippedCount := 0

	for _, t := range templates {
		outputPath := filepath.Join(t.outputPath, t.fileName)
		relPath := FormatRelativePath(cwd, outputPath)

		if engine.FileExists(outputPath) {
			if !jsonOutput {
				log.Warning("File exists, skipping", "file", relPath)
			}
			tracker.AddSkipped(relPath)
			skippedCount++
			continue
		}

		if err := engine.RenderAndWrite(t.template, outputPath, data); err != nil {
			if jsonOutput {
				tracker.AddError(fmt.Sprintf("failed to generate %s: %s", t.fileName, err.Error()))
				continue
			}
			return generatedCount, skippedCount, fmt.Errorf("failed to generate %s: %w", t.fileName, err)
		}

		if !jsonOutput {
			log.Info("Created", "file", relPath)
		}
		tracker.AddGenerated(relPath)
		generatedCount++
	}

	return generatedCount, skippedCount, nil
}

func getSecurityPackage(profile *detector.ProjectProfile) string {
	switch profile.Architecture {
	case detector.ArchFeature:
		return profile.BasePackage + ".common.security"
	case detector.ArchHexagonal:
		return profile.BasePackage + ".infrastructure.security"
	case detector.ArchClean:
		return profile.BasePackage + ".infrastructure.security"
	default:
		return profile.BasePackage + ".security"
	}
}

func getUserEntityPackage(profile *detector.ProjectProfile) string {
	switch profile.Architecture {
	case detector.ArchFeature:
		return profile.BasePackage + ".user"
	case detector.ArchHexagonal:
		return profile.BasePackage + ".domain.model"
	case detector.ArchClean:
		return profile.BasePackage + ".domain.entity"
	default:
		return profile.BasePackage + ".entity"
	}
}

func getUserRepositoryPackage(profile *detector.ProjectProfile) string {
	switch profile.Architecture {
	case detector.ArchFeature:
		return profile.BasePackage + ".user"
	case detector.ArchHexagonal:
		return profile.BasePackage + ".adapter.persistence"
	case detector.ArchClean:
		return profile.BasePackage + ".infrastructure.persistence"
	default:
		return profile.BasePackage + ".repository"
	}
}

func buildSecurityTemplateData(profile *detector.ProjectProfile, securityPackage, userEntityPackage, userRepositoryPackage string, cfg securityConfig) map[string]any {
	validationImport := "jakarta.validation"
	if profile.ValidationStyle == detector.ValidationJavax {
		validationImport = "javax.validation"
	}

	idType := "Long"
	idStrategy := "IDENTITY"
	if profile.IDType == "UUID" {
		idType = "UUID"
		idStrategy = "UUID"
	}

	return map[string]any{
		"BasePackage":           profile.BasePackage,
		"SecurityPackage":       securityPackage,
		"UserEntityPackage":     userEntityPackage,
		"UserRepositoryPackage": userRepositoryPackage,
		"UserEntityName":        cfg.UserEntityName,
		"HasLombok":             profile.Lombok.Detected,
		"HasValidation":         profile.HasValidation,
		"ValidationImport":      validationImport,
		"Architecture":          string(profile.Architecture),
		"IdType":                idType,
		"IdStrategy":            idStrategy,
		"DefaultJwtSecret":      generateJwtSecret(),
		"GenerateEntities":      cfg.GenerateEntities,
	}
}

func generateJwtSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "your-256-bit-secret-key-here-change-in-production"
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func printSecurityInstructions(types []SecurityType) {
	log := logger.Default()
	log.Info("")
	log.Info("Next steps:")

	for _, t := range types {
		switch t {
		case SecurityJWT:
			log.Info("  JWT Authentication:")
			log.Info("    1. Add jwt.secret to application.properties/yml")
			log.Info("    2. Configure jwt.expiration (default: 86400000ms = 24h)")
			log.Info("    3. Use POST /api/auth/register to create users")
			log.Info("    4. Use POST /api/auth/login to get tokens")
		case SecuritySession:
			log.Info("  Session Authentication:")
			log.Info("    1. Create remember-me token table:")
			log.Info("       CREATE TABLE persistent_logins (")
			log.Info("         username VARCHAR(64) NOT NULL,")
			log.Info("         series VARCHAR(64) PRIMARY KEY,")
			log.Info("         token VARCHAR(64) NOT NULL,")
			log.Info("         last_used TIMESTAMP NOT NULL")
			log.Info("       );")
			log.Info("    2. Create login.html and register.html templates")
		case SecurityOAuth2:
			log.Info("  OAuth2 Authentication:")
			log.Info("    1. Add OAuth2 client configuration to application.yml:")
			log.Info("       spring.security.oauth2.client.registration.google.client-id=xxx")
			log.Info("       spring.security.oauth2.client.registration.google.client-secret=xxx")
		}
	}
}
