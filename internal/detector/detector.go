package detector

import (
	"time"

	"github.com/spf13/afero"
)

const (
	DefaultCacheMaxAge = 24 * time.Hour
)

type Detector struct {
	fs                   afero.Fs
	projectDir           string
	scanner              *Scanner
	confidenceCalculator *ConfidenceCalculator
	cacheMaxAge          time.Duration
}

type DetectorOption func(*Detector)

func WithCacheMaxAge(maxAge time.Duration) DetectorOption {
	return func(d *Detector) {
		d.cacheMaxAge = maxAge
	}
}

func WithFileSystem(fs afero.Fs) DetectorOption {
	return func(d *Detector) {
		d.fs = fs
		d.scanner = NewScanner(fs, d.projectDir)
	}
}

func NewDetector(projectDir string, opts ...DetectorOption) *Detector {
	fs := afero.NewOsFs()
	d := &Detector{
		fs:                   fs,
		projectDir:           projectDir,
		scanner:              NewScanner(fs, projectDir),
		confidenceCalculator: NewConfidenceCalculator(),
		cacheMaxAge:          DefaultCacheMaxAge,
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

func (d *Detector) Detect() (*ProjectProfile, error) {
	scanResult, err := d.scanner.Scan()
	if err != nil {
		return nil, err
	}

	profile := NewEmptyProfile()
	profile.ProjectRoot = d.projectDir
	profile.BasePackage = scanResult.BasePackage
	profile.SourceRoot = scanResult.SourceRoot
	profile.TestRoot = scanResult.TestRoot

	d.detectArchitecture(scanResult, profile)
	d.detectFeatureStyle(scanResult, profile)
	d.detectBaseClasses(scanResult, profile)
	d.detectNamingConventions(scanResult, profile)
	d.detectIDType(scanResult, profile)
	d.detectMapper(scanResult, profile)
	d.detectLombok(scanResult, profile)
	d.detectExceptions(scanResult, profile)
	d.detectSwagger(scanResult, profile)
	d.detectValidation(scanResult, profile)
	d.detectDatabase(scanResult, profile)
	d.detectTestingProfile(scanResult, profile)

	return profile, nil
}

func (d *Detector) detectArchitecture(scan *ScanResult, profile *ProjectProfile) {
	if len(scan.SourceFiles) == 0 {
		profile.Architecture = ArchLayered
		profile.ArchConfidence = 1.0
		return
	}

	scores := make(map[ArchitectureType]float64)

	scores[ArchLayered] = d.calculateLayeredScore(scan)
	scores[ArchFeature] = d.calculateFeatureScore(scan)
	scores[ArchHexagonal] = d.calculateHexagonalScore(scan)
	scores[ArchClean] = d.calculateCleanScore(scan)
	scores[ArchModular] = d.calculateModularScore(scan)
	scores[ArchFlat] = d.calculateFlatScore(scan)

	bestArch := ArchLayered
	bestScore := scores[ArchLayered]

	for arch, score := range scores {
		if score > bestScore {
			bestScore = score
			bestArch = arch
		}
	}

	profile.Architecture = bestArch
	profile.ArchConfidence = bestScore

	if bestArch == ArchFeature {
		profile.FeatureModules = d.extractFeatureModules(scan, profile.BasePackage)
	}
}

func (d *Detector) calculateLayeredScore(scan *ScanResult) float64 {
	layerPackages := map[string]bool{
		"controller": false,
		"service":    false,
		"repository": false,
		"entity":     false,
		"dto":        false,
		"model":      false,
	}

	filesInLayers := 0

	for _, file := range scan.SourceFiles {
		for layer := range layerPackages {
			if containsPackagePart(file.Package, layer) {
				if isDirectLayerPackage(file.Package, scan.BasePackage, layer) {
					layerPackages[layer] = true
					filesInLayers++
				}
			}
		}
	}

	foundLayers := 0
	for _, found := range layerPackages {
		if found {
			foundLayers++
		}
	}

	if len(scan.SourceFiles) == 0 {
		return 0.0
	}

	layerCoverage := float64(foundLayers) / 4.0
	fileRatio := float64(filesInLayers) / float64(len(scan.SourceFiles))

	return d.confidenceCalculator.Calculate(
		(layerCoverage+fileRatio)/2,
		len(scan.SourceFiles),
		layerCoverage,
	)
}

func (d *Detector) calculateFeatureScore(scan *ScanResult) float64 {
	if scan.BasePackage == "" {
		return 0.0
	}

	featureModules := d.extractFeatureModules(scan, scan.BasePackage)

	if len(featureModules) < 2 {
		return 0.0
	}

	modulesWithMultipleTypes := 0
	totalFilesInModules := 0

	for _, module := range featureModules {
		typesInModule := d.countFileTypesInModule(scan, scan.BasePackage, module)
		if typesInModule >= 2 {
			modulesWithMultipleTypes++
		}
		totalFilesInModules += d.countFilesInModule(scan, scan.BasePackage, module)
	}

	hasCommonPackage := d.hasCommonPackage(scan, scan.BasePackage)

	moduleQuality := float64(modulesWithMultipleTypes) / float64(len(featureModules))
	fileRatio := float64(totalFilesInModules) / float64(len(scan.SourceFiles))

	score := d.confidenceCalculator.CalculateArchitectureConfidence(
		totalFilesInModules,
		len(scan.SourceFiles),
		hasCommonPackage,
		1.0-moduleQuality,
	)

	if hasCommonPackage {
		score += 0.05
	}
	if moduleQuality > 0.7 {
		score += 0.05
	}

	return clamp(score*fileRatio+score*(1-fileRatio)*0.5, 0.0, 1.0)
}

func (d *Detector) calculateHexagonalScore(scan *ScanResult) float64 {
	hexagonalMarkers := []string{"domain", "application", "infrastructure", "adapter", "port"}
	foundMarkers := 0

	packages := d.scanner.GetUniquePackages(scan.SourceFiles)

	for _, marker := range hexagonalMarkers {
		for _, pkg := range packages {
			if containsPackagePart(pkg, marker) {
				foundMarkers++
				break
			}
		}
	}

	if foundMarkers < 2 {
		return 0.0
	}

	return d.confidenceCalculator.Calculate(
		float64(foundMarkers)/float64(len(hexagonalMarkers)),
		len(scan.SourceFiles),
		float64(foundMarkers)/float64(len(hexagonalMarkers)),
	)
}

func (d *Detector) calculateCleanScore(scan *ScanResult) float64 {
	cleanSpecificMarkers := []string{"usecase", "gateway", "presenter", "interactor"}
	cleanSupportMarkers := []string{"domain", "application", "infrastructure"}
	foundSpecific := 0
	foundSupport := 0

	packages := d.scanner.GetUniquePackages(scan.SourceFiles)

	for _, marker := range cleanSpecificMarkers {
		for _, pkg := range packages {
			if containsPackagePart(pkg, marker) {
				foundSpecific++
				break
			}
		}
	}

	for _, marker := range cleanSupportMarkers {
		for _, pkg := range packages {
			if containsPackagePart(pkg, marker) {
				foundSupport++
				break
			}
		}
	}

	if foundSpecific < 1 {
		return 0.0
	}

	specificRatio := float64(foundSpecific) / float64(len(cleanSpecificMarkers))
	supportRatio := float64(foundSupport) / float64(len(cleanSupportMarkers))
	combinedRatio := (specificRatio*0.7 + supportRatio*0.3)

	baseScore := d.confidenceCalculator.Calculate(
		combinedRatio,
		len(scan.SourceFiles),
		combinedRatio,
	)

	if foundSpecific >= 2 {
		baseScore += 0.15
	}
	if foundSupport >= 2 {
		baseScore += 0.1
	}

	return clamp(baseScore, 0.0, 1.0)
}

func (d *Detector) calculateModularScore(scan *ScanResult) float64 {
	modularMarkers := []string{"api", "internal", "module"}
	foundMarkers := 0

	packages := d.scanner.GetUniquePackages(scan.SourceFiles)

	for _, marker := range modularMarkers {
		for _, pkg := range packages {
			if containsPackagePart(pkg, marker) {
				foundMarkers++
				break
			}
		}
	}

	if foundMarkers < 2 {
		return 0.0
	}

	return d.confidenceCalculator.Calculate(
		float64(foundMarkers)/float64(len(modularMarkers)),
		len(scan.SourceFiles),
		0.7,
	)
}

func (d *Detector) calculateFlatScore(scan *ScanResult) float64 {
	if len(scan.SourceFiles) > 15 {
		return 0.0
	}

	uniquePackages := d.scanner.GetUniquePackages(scan.SourceFiles)

	if len(uniquePackages) <= 2 {
		return d.confidenceCalculator.Calculate(
			1.0,
			len(scan.SourceFiles),
			1.0,
		)
	}

	return 0.0
}

func (d *Detector) detectFeatureStyle(scan *ScanResult, profile *ProjectProfile) {
	if profile.Architecture != ArchFeature {
		profile.FeatureStyle = ""
		return
	}

	if len(profile.FeatureModules) == 0 {
		profile.FeatureStyle = FeatureStyleNested
		return
	}

	flatCount := 0
	nestedCount := 0

	for _, module := range profile.FeatureModules {
		modulePrefix := profile.BasePackage + "." + module

		for _, file := range scan.SourceFiles {
			if !hasPackagePrefix(file.Package, modulePrefix) {
				continue
			}

			if file.Package == modulePrefix {
				flatCount++
			} else {
				nestedCount++
			}
		}
	}

	if flatCount > nestedCount {
		profile.FeatureStyle = FeatureStyleFlat
	} else {
		profile.FeatureStyle = FeatureStyleNested
	}
}

func (d *Detector) extractFeatureModules(scan *ScanResult, basePackage string) []string {
	if basePackage == "" {
		return nil
	}

	moduleSet := make(map[string]bool)
	baseParts := len(splitPackage(basePackage))

	for _, file := range scan.SourceFiles {
		parts := splitPackage(file.Package)
		if len(parts) > baseParts {
			moduleName := parts[baseParts]
			if !isLayerName(moduleName) && moduleName != "config" {
				moduleSet[moduleName] = true
			}
		}
	}

	var modules []string
	for module := range moduleSet {
		modules = append(modules, module)
	}

	return modules
}

func (d *Detector) countFileTypesInModule(scan *ScanResult, basePackage, module string) int {
	typeSet := make(map[JavaFileType]bool)
	modulePrefix := basePackage + "." + module

	for _, file := range scan.SourceFiles {
		if hasPackagePrefix(file.Package, modulePrefix) {
			typeSet[file.FileType] = true
		}
	}

	return len(typeSet)
}

func (d *Detector) countFilesInModule(scan *ScanResult, basePackage, module string) int {
	count := 0
	modulePrefix := basePackage + "." + module

	for _, file := range scan.SourceFiles {
		if hasPackagePrefix(file.Package, modulePrefix) {
			count++
		}
	}

	return count
}

func (d *Detector) hasCommonPackage(scan *ScanResult, basePackage string) bool {
	commonPrefix := basePackage + ".common"

	for _, file := range scan.SourceFiles {
		if hasPackagePrefix(file.Package, commonPrefix) {
			return true
		}
	}

	return false
}

func (d *Detector) detectBaseClasses(scan *ScanResult, profile *ProjectProfile) {
	entities := d.scanner.GetFilesByType(scan.SourceFiles, FileTypeEntity)

	if len(entities) == 0 {
		return
	}

	parentCounts := make(map[string]int)
	parentFiles := make(map[string]*JavaFile)

	for _, entity := range entities {
		if entity.ExtendsClass != "" && entity.ExtendsClass != "Object" {
			parentCounts[entity.ExtendsClass]++
		}
	}

	for _, file := range scan.SourceFiles {
		if _, exists := parentCounts[file.ClassName]; exists {
			parentFiles[file.ClassName] = file
		}
	}

	var bestParent string
	var bestCount int

	for parent, count := range parentCounts {
		if count > bestCount {
			bestCount = count
			bestParent = parent
		}
	}

	if bestCount >= 2 || (bestCount == 1 && len(entities) <= 3) {
		if parentFile, exists := parentFiles[bestParent]; exists {
			profile.BaseEntity = &BaseClassInfo{
				Name:     bestParent,
				Package:  parentFile.Package,
				FullPath: parentFile.Path,
			}
		} else {
			profile.BaseEntity = &BaseClassInfo{
				Name: bestParent,
			}
		}
	}
}

func (d *Detector) detectNamingConventions(scan *ScanResult, profile *ProjectProfile) {
	d.detectDTONaming(scan, profile)
	d.detectControllerSuffix(scan, profile)
}

func (d *Detector) detectDTONaming(scan *ScanResult, profile *ProjectProfile) {
	dtos := d.scanner.GetFilesByType(scan.SourceFiles, FileTypeDTO)

	requestResponseCount := 0
	dtoUpperCount := 0
	dtoLowerCount := 0

	for _, dto := range dtos {
		name := dto.ClassName
		switch {
		case endsWithAny(name, "Request", "Response"):
			requestResponseCount++
		case endsWith(name, "DTO"):
			dtoUpperCount++
		case endsWith(name, "Dto"):
			dtoLowerCount++
		}
	}

	if requestResponseCount >= dtoUpperCount && requestResponseCount >= dtoLowerCount {
		profile.DTONaming = DTONamingRequestResponse
	} else if dtoUpperCount >= dtoLowerCount {
		profile.DTONaming = DTONamingDTOUpper
	} else {
		profile.DTONaming = DTONamingDTOLower
	}
}

func (d *Detector) detectControllerSuffix(scan *ScanResult, profile *ProjectProfile) {
	controllers := d.scanner.GetFilesByType(scan.SourceFiles, FileTypeController)
	resourceCount := 0
	controllerCount := 0

	for _, c := range controllers {
		if endsWith(c.ClassName, "Resource") {
			resourceCount++
		} else if endsWith(c.ClassName, "Controller") {
			controllerCount++
		}
	}

	if resourceCount > controllerCount {
		profile.ControllerSuffix = "Resource"
	} else {
		profile.ControllerSuffix = "Controller"
	}
}

func (d *Detector) detectIDType(scan *ScanResult, profile *ProjectProfile) {
	entities := d.scanner.GetFilesByType(scan.SourceFiles, FileTypeEntity)

	if len(entities) == 0 {
		profile.IDType = "Long"
		return
	}

	idTypeCounts := make(map[string]int)

	for _, entity := range entities {
		for _, imp := range entity.Imports {
			switch {
			case endsWith(imp, "UUID"):
				idTypeCounts["UUID"]++
			}
		}
	}

	if idTypeCounts["UUID"] > 0 {
		profile.IDType = "UUID"
		profile.IDAnnotation = "@GeneratedValue(strategy = GenerationType.UUID)"
	} else {
		profile.IDType = "Long"
		profile.IDAnnotation = "@GeneratedValue(strategy = GenerationType.IDENTITY)"
	}
}

func (d *Detector) detectMapper(scan *ScanResult, profile *ProjectProfile) {
	mappers := d.scanner.GetFilesByAnnotation(scan.SourceFiles, "Mapper")

	if len(mappers) > 0 {
		profile.Mapper = MapperMapStruct
		return
	}

	for _, file := range scan.SourceFiles {
		for _, imp := range file.Imports {
			if containsString(imp, "mapstruct") {
				profile.Mapper = MapperMapStruct
				return
			}
			if containsString(imp, "modelmapper") {
				profile.Mapper = MapperModelMapper
				return
			}
		}
	}

	profile.Mapper = MapperManual
}

func (d *Detector) detectLombok(scan *ScanResult, profile *ProjectProfile) {
	lombokAnnotations := map[string]*bool{
		"Data":                    &profile.Lombok.UseData,
		"Builder":                 &profile.Lombok.UseBuilder,
		"Getter":                  &profile.Lombok.UseAccessors,
		"Setter":                  &profile.Lombok.UseAccessors,
		"Slf4j":                   &profile.Lombok.UseSlf4j,
		"RequiredArgsConstructor": &profile.Lombok.UseRequiredArgs,
		"AllArgsConstructor":      &profile.Lombok.UseAllArgs,
		"NoArgsConstructor":       &profile.Lombok.UseNoArgs,
	}

	for _, file := range scan.SourceFiles {
		for _, ann := range file.Annotations {
			if ptr, exists := lombokAnnotations[ann]; exists {
				*ptr = true
				profile.Lombok.Detected = true
			}
		}
	}
}

func (d *Detector) detectExceptions(scan *ScanResult, profile *ProjectProfile) {
	exceptionFiles := d.scanner.GetFilesByType(scan.SourceFiles, FileTypeException)

	for _, file := range exceptionFiles {
		for _, ann := range file.Annotations {
			if ann == "RestControllerAdvice" || ann == "ControllerAdvice" {
				profile.Exceptions.HasGlobalHandler = true
				profile.Exceptions.HandlerPackage = file.Package
			}
		}

		if endsWith(file.ClassName, "Exception") {
			profile.Exceptions.CustomExceptions = append(
				profile.Exceptions.CustomExceptions,
				ExceptionInfo{
					Name:    file.ClassName,
					Package: file.Package,
				},
			)
		}
	}
}

func (d *Detector) detectSwagger(scan *ScanResult, profile *ProjectProfile) {
	for _, file := range scan.SourceFiles {
		for _, imp := range file.Imports {
			if containsString(imp, "io.swagger.v3") || containsString(imp, "springdoc") {
				profile.HasSwagger = true
				profile.SwaggerStyle = SwaggerOpenAPI3
				return
			}
			if containsString(imp, "io.swagger") && !containsString(imp, "v3") {
				profile.HasSwagger = true
				profile.SwaggerStyle = SwaggerV2
				return
			}
		}

		for _, ann := range file.Annotations {
			if ann == "Operation" || ann == "Tag" || ann == "ApiResponse" {
				profile.HasSwagger = true
				profile.SwaggerStyle = SwaggerOpenAPI3
				return
			}
			if ann == "Api" || ann == "ApiOperation" {
				profile.HasSwagger = true
				profile.SwaggerStyle = SwaggerV2
				return
			}
		}
	}
}

func (d *Detector) detectValidation(scan *ScanResult, profile *ProjectProfile) {
	for _, file := range scan.SourceFiles {
		for _, imp := range file.Imports {
			if containsString(imp, "jakarta.validation") {
				profile.HasValidation = true
				profile.ValidationStyle = ValidationJakarta
				return
			}
			if containsString(imp, "javax.validation") {
				profile.HasValidation = true
				profile.ValidationStyle = ValidationJavax
				return
			}
		}

		for _, ann := range file.Annotations {
			if ann == "Valid" || ann == "NotNull" || ann == "NotBlank" || ann == "Size" {
				profile.HasValidation = true
				profile.ValidationStyle = ValidationJakarta
				return
			}
		}
	}
}

func (d *Detector) detectDatabase(scan *ScanResult, profile *ProjectProfile) {
	hasJPA := false
	hasCassandra := false
	hasMongo := false
	hasR2DBC := false

	for _, file := range scan.SourceFiles {
		for _, ann := range file.Annotations {
			switch ann {
			case "Entity", "Table":
				hasJPA = true
			case "Document":
				hasMongo = true
			}
		}

		for _, imp := range file.Imports {
			if containsString(imp, "cassandra") {
				hasCassandra = true
			}
			if containsString(imp, "r2dbc") {
				hasR2DBC = true
			}
			if containsString(imp, "mongodb") {
				hasMongo = true
			}
		}

		for _, iface := range file.ImplementsInterfaces {
			if containsString(iface, "CassandraRepository") {
				hasCassandra = true
			}
			if containsString(iface, "MongoRepository") {
				hasMongo = true
			}
			if containsString(iface, "R2dbcRepository") || containsString(iface, "ReactiveCrudRepository") {
				hasR2DBC = true
			}
		}
	}

	dbCount := 0
	if hasJPA {
		dbCount++
	}
	if hasCassandra {
		dbCount++
	}
	if hasMongo {
		dbCount++
	}
	if hasR2DBC {
		dbCount++
	}

	if dbCount > 1 {
		profile.Database = DatabaseMulti
	} else if hasR2DBC {
		profile.Database = DatabaseR2DBC
	} else if hasMongo {
		profile.Database = DatabaseMongo
	} else if hasCassandra {
		profile.Database = DatabaseCassandra
	} else if hasJPA {
		profile.Database = DatabaseJPA
	} else {
		profile.Database = DatabaseJPA
	}
}

func (d *Detector) detectTestingProfile(scan *ScanResult, profile *ProjectProfile) {
	profile.Testing.Framework = "junit5"

	for _, file := range scan.TestFiles {
		for _, imp := range file.Imports {
			if containsString(imp, "mockito") {
				profile.Testing.HasMockito = true
			}
			if containsString(imp, "testcontainers") {
				profile.Testing.HasTestcontainers = true
			}
			if containsString(imp, "rest-assured") || containsString(imp, "restassured") {
				profile.Testing.HasRestAssured = true
			}
		}

		for _, ann := range file.Annotations {
			if ann == "Testcontainers" {
				profile.Testing.HasTestcontainers = true
			}
		}
	}

	profile.Testing.StructureMirror = d.detectTestStructureMirror(scan)
}

func (d *Detector) detectTestStructureMirror(scan *ScanResult) bool {
	if len(scan.TestFiles) == 0 || len(scan.SourceFiles) == 0 {
		return true
	}

	sourcePackages := make(map[string]bool)
	for _, f := range scan.SourceFiles {
		sourcePackages[f.Package] = true
	}

	matchingPackages := 0
	for _, f := range scan.TestFiles {
		if sourcePackages[f.Package] {
			matchingPackages++
		}
	}

	return float64(matchingPackages)/float64(len(scan.TestFiles)) > 0.5
}

func (d *Detector) GetScanner() *Scanner {
	return d.scanner
}

func (d *Detector) GetConfidenceCalculator() *ConfidenceCalculator {
	return d.confidenceCalculator
}

func containsPackagePart(pkg, part string) bool {
	parts := splitPackage(pkg)
	for _, p := range parts {
		if p == part {
			return true
		}
	}
	return false
}

func isDirectLayerPackage(pkg, basePackage, layer string) bool {
	expected := basePackage + "." + layer
	return pkg == expected || hasPackagePrefix(pkg, expected+".")
}

func hasPackagePrefix(pkg, prefix string) bool {
	return pkg == prefix || (len(pkg) > len(prefix) && pkg[:len(prefix)] == prefix && pkg[len(prefix)] == '.')
}

func splitPackage(pkg string) []string {
	if pkg == "" {
		return nil
	}

	var parts []string
	start := 0
	for i, c := range pkg {
		if c == '.' {
			parts = append(parts, pkg[start:i])
			start = i + 1
		}
	}
	parts = append(parts, pkg[start:])
	return parts
}

func isLayerName(name string) bool {
	layers := []string{
		"controller", "service", "repository", "entity", "dto", "mapper", "model", "exception", "config",
		"domain", "application", "infrastructure", "adapter", "port",
		"usecase", "gateway", "presenter", "interactor",
		"api", "internal", "module", "web", "persistence", "common",
	}
	for _, layer := range layers {
		if name == layer {
			return true
		}
	}
	return false
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func endsWithAny(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if endsWith(s, suffix) {
			return true
		}
	}
	return false
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
