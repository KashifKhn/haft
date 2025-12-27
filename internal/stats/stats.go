package stats

import (
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/boyter/scc/v3/processor"
)

type LanguageStats struct {
	Name       string
	Files      int64
	Lines      int64
	Code       int64
	Comments   int64
	Blanks     int64
	Complexity int64
	Bytes      int64
}

type ProjectStats struct {
	Languages       []LanguageStats
	TotalFiles      int64
	TotalLines      int64
	TotalCode       int64
	TotalComments   int64
	TotalBlanks     int64
	TotalComplexity int64
	TotalBytes      int64
	EstimatedCost   float64
	EstimatedMonths float64
	EstimatedPeople float64
}

var initOnce sync.Once

func initSCC() {
	initOnce.Do(func() {
		processor.ProcessConstants()
	})
}

func CountProject(dir string) (*ProjectStats, error) {
	initSCC()

	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	languageMap := make(map[string]*LanguageStats)

	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			base := filepath.Base(path)
			if base == ".git" || base == ".svn" || base == ".hg" || base == "node_modules" || base == "target" || base == "build" || base == ".gradle" || base == ".idea" {
				return filepath.SkipDir
			}
			return nil
		}

		if info.Size() == 0 {
			return nil
		}

		possibleLanguages, _ := processor.DetectLanguage(path)
		if len(possibleLanguages) == 0 {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		language := processor.DetermineLanguage(filepath.Base(path), "", possibleLanguages, content)
		if language == "" {
			return nil
		}

		fileJob := &processor.FileJob{
			Filename:          filepath.Base(path),
			Location:          path,
			Language:          language,
			PossibleLanguages: possibleLanguages,
			Content:           content,
			Bytes:             int64(len(content)),
		}

		processor.CountStats(fileJob)

		if fileJob.Binary {
			return nil
		}

		stats, exists := languageMap[language]
		if !exists {
			stats = &LanguageStats{Name: language}
			languageMap[language] = stats
		}

		stats.Files++
		stats.Lines += fileJob.Lines
		stats.Code += fileJob.Code
		stats.Comments += fileJob.Comment
		stats.Blanks += fileJob.Blank
		stats.Complexity += fileJob.Complexity
		stats.Bytes += fileJob.Bytes

		return nil
	})

	if err != nil {
		return nil, err
	}

	result := &ProjectStats{
		Languages: make([]LanguageStats, 0, len(languageMap)),
	}

	for _, stats := range languageMap {
		result.Languages = append(result.Languages, *stats)
		result.TotalFiles += stats.Files
		result.TotalLines += stats.Lines
		result.TotalCode += stats.Code
		result.TotalComments += stats.Comments
		result.TotalBlanks += stats.Blanks
		result.TotalComplexity += stats.Complexity
		result.TotalBytes += stats.Bytes
	}

	sort.Slice(result.Languages, func(i, j int) bool {
		return result.Languages[i].Code > result.Languages[j].Code
	})

	if result.TotalCode > 0 {
		effort := processor.EstimateEffort(result.TotalCode, processor.EAF)
		result.EstimatedCost = processor.EstimateCost(effort, processor.AverageWage, processor.Overhead)
		result.EstimatedMonths = processor.EstimateScheduleMonths(effort)
		if result.EstimatedMonths > 0 {
			result.EstimatedPeople = effort / result.EstimatedMonths
		}
	}

	return result, nil
}

func CountProjectQuick(dir string) (*ProjectStats, error) {
	initSCC()

	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	var totalLines, totalCode, totalComments, totalBlanks, totalFiles int64

	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			base := filepath.Base(path)
			if base == ".git" || base == ".svn" || base == ".hg" || base == "node_modules" || base == "target" || base == "build" || base == ".gradle" || base == ".idea" {
				return filepath.SkipDir
			}
			return nil
		}

		if info.Size() == 0 {
			return nil
		}

		possibleLanguages, _ := processor.DetectLanguage(path)
		if len(possibleLanguages) == 0 {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		language := processor.DetermineLanguage(filepath.Base(path), "", possibleLanguages, content)
		if language == "" {
			return nil
		}

		fileJob := &processor.FileJob{
			Filename:          filepath.Base(path),
			Location:          path,
			Language:          language,
			PossibleLanguages: possibleLanguages,
			Content:           content,
			Bytes:             int64(len(content)),
		}

		processor.CountStats(fileJob)

		if fileJob.Binary {
			return nil
		}

		totalFiles++
		totalLines += fileJob.Lines
		totalCode += fileJob.Code
		totalComments += fileJob.Comment
		totalBlanks += fileJob.Blank

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ProjectStats{
		TotalFiles:    totalFiles,
		TotalLines:    totalLines,
		TotalCode:     totalCode,
		TotalComments: totalComments,
		TotalBlanks:   totalBlanks,
	}, nil
}
