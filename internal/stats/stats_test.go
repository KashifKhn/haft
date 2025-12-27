package stats

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountProject(t *testing.T) {
	tmpDir := t.TempDir()

	javaFile := filepath.Join(tmpDir, "Test.java")
	javaContent := `package com.example;

public class Test {
    // This is a comment
    public static void main(String[] args) {
        System.out.println("Hello");
    }
}
`
	err := os.WriteFile(javaFile, []byte(javaContent), 0644)
	require.NoError(t, err)

	propsFile := filepath.Join(tmpDir, "application.properties")
	propsContent := `# Database config
spring.datasource.url=jdbc:mysql://localhost:3306/test
spring.datasource.username=root
`
	err = os.WriteFile(propsFile, []byte(propsContent), 0644)
	require.NoError(t, err)

	stats, err := CountProject(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, stats)

	assert.GreaterOrEqual(t, stats.TotalFiles, int64(1))
	assert.GreaterOrEqual(t, stats.TotalLines, int64(1))
	assert.GreaterOrEqual(t, stats.TotalCode, int64(1))
}

func TestCountProjectQuick(t *testing.T) {
	tmpDir := t.TempDir()

	javaFile := filepath.Join(tmpDir, "Test.java")
	javaContent := `package com.example;

public class Test {
    // Comment
    public void run() {}
}
`
	err := os.WriteFile(javaFile, []byte(javaContent), 0644)
	require.NoError(t, err)

	stats, err := CountProjectQuick(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, stats)

	assert.GreaterOrEqual(t, stats.TotalFiles, int64(1))
	assert.GreaterOrEqual(t, stats.TotalLines, int64(1))
}

func TestCountProject_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	stats, err := CountProject(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(0), stats.TotalFiles)
}

func TestCountProject_SkipsIgnoredDirs(t *testing.T) {
	tmpDir := t.TempDir()

	gitDir := filepath.Join(tmpDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	require.NoError(t, err)

	gitFile := filepath.Join(gitDir, "config")
	err = os.WriteFile(gitFile, []byte("gitconfig"), 0644)
	require.NoError(t, err)

	targetDir := filepath.Join(tmpDir, "target")
	err = os.MkdirAll(targetDir, 0755)
	require.NoError(t, err)

	targetFile := filepath.Join(targetDir, "Test.class")
	err = os.WriteFile(targetFile, []byte("classfile"), 0644)
	require.NoError(t, err)

	srcDir := filepath.Join(tmpDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "Main.java")
	javaContent := `public class Main {
    public static void main(String[] args) {}
}
`
	err = os.WriteFile(srcFile, []byte(javaContent), 0644)
	require.NoError(t, err)

	stats, err := CountProject(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, stats)

	assert.Equal(t, int64(1), stats.TotalFiles)
}

func TestCountProject_CurrentDirDefault(t *testing.T) {
	stats, err := CountProject("")
	require.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestCountProjectQuick_CurrentDirDefault(t *testing.T) {
	stats, err := CountProjectQuick("")
	require.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestCountProject_LanguagesSorted(t *testing.T) {
	tmpDir := t.TempDir()

	javaFile := filepath.Join(tmpDir, "Main.java")
	javaContent := `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello");
        System.out.println("World");
        System.out.println("Java");
    }
}
`
	err := os.WriteFile(javaFile, []byte(javaContent), 0644)
	require.NoError(t, err)

	propsFile := filepath.Join(tmpDir, "app.properties")
	propsContent := `key=value
`
	err = os.WriteFile(propsFile, []byte(propsContent), 0644)
	require.NoError(t, err)

	stats, err := CountProject(tmpDir)
	require.NoError(t, err)

	if len(stats.Languages) >= 2 {
		assert.GreaterOrEqual(t, stats.Languages[0].Code, stats.Languages[1].Code)
	}
}

func TestCountProject_COCOMOEstimates(t *testing.T) {
	tmpDir := t.TempDir()

	javaFile := filepath.Join(tmpDir, "Main.java")
	javaContent := `public class Main {
    public static void main(String[] args) {
        System.out.println("Line 1");
        System.out.println("Line 2");
        System.out.println("Line 3");
    }
}
`
	err := os.WriteFile(javaFile, []byte(javaContent), 0644)
	require.NoError(t, err)

	stats, err := CountProject(tmpDir)
	require.NoError(t, err)

	if stats.TotalCode > 0 {
		assert.Greater(t, stats.EstimatedCost, float64(0))
		assert.Greater(t, stats.EstimatedMonths, float64(0))
	}
}
