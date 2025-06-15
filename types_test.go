package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test String methods for custom types
func TestLayerName_String(t *testing.T) {
	ln := LayerName("test-layer")
	expected := "test-layer"
	assert.Equal(t, expected, ln.String())
}

func TestLayerPath_String(t *testing.T) {
	lp := LayerPath("path/to/layer")
	expected := "path/to/layer"
	assert.Equal(t, expected, lp.String())
}

func TestModuleName_String(t *testing.T) {
	mn := ModuleName("github.com/example/module")
	expected := "github.com/example/module"
	assert.Equal(t, expected, mn.String())
}

func TestPackageName_String(t *testing.T) {
	pn := PackageName("main")
	expected := "main"
	assert.Equal(t, expected, pn.String())
}

func TestFilePath_String(t *testing.T) {
	fp := FilePath("path/to/file.go")
	expected := "path/to/file.go"
	assert.Equal(t, expected, fp.String())
}

// Test LayerName validation
func TestLayerName_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    LayerName
		expected bool
	}{
		{"valid layer name", LayerName("valid-layer"), true},
		{"empty string", LayerName(""), false},
		{"whitespace only", LayerName("   "), false},
		{"single character", LayerName("a"), true},
		{"with spaces", LayerName("layer with spaces"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestLayerName_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     LayerName
		expectErr bool
	}{
		{"valid layer name", LayerName("valid-layer"), false},
		{"empty string", LayerName(""), true},
		{"whitespace only", LayerName("   "), true},
		{"single character", LayerName("a"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test LayerPath validation
func TestLayerPath_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    LayerPath
		expected bool
	}{
		{"valid path", LayerPath("path/to/layer"), true},
		{"empty string", LayerPath(""), false},
		{"whitespace only", LayerPath("   "), false},
		{"path with parent directory", LayerPath("path/../other"), false},
		{"path with double dots", LayerPath("path/.."), false},
		{"single dot", LayerPath("."), true},
		{"root path", LayerPath("/"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestLayerPath_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     LayerPath
		expectErr bool
	}{
		{"valid path", LayerPath("path/to/layer"), false},
		{"empty string", LayerPath(""), true},
		{"whitespace only", LayerPath("   "), true},
		{"path with parent directory", LayerPath("path/../other"), true},
		{"path with double dots", LayerPath("path/.."), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test ModuleName validation
func TestModuleName_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    ModuleName
		expected bool
	}{
		{"valid module name", ModuleName("github.com/example/module"), true},
		{"empty string", ModuleName(""), false},
		{"whitespace only", ModuleName("   "), false},
		{"module with spaces", ModuleName("module with spaces"), false},
		{"single word", ModuleName("module"), true},
		{"with dashes", ModuleName("my-module"), true},
		{"with underscores", ModuleName("my_module"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestModuleName_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     ModuleName
		expectErr bool
	}{
		{"valid module name", ModuleName("github.com/example/module"), false},
		{"empty string", ModuleName(""), true},
		{"whitespace only", ModuleName("   "), true},
		{"module with spaces", ModuleName("module with spaces"), true},
		{"single word", ModuleName("module"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test PackageName validation
func TestPackageName_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    PackageName
		expected bool
	}{
		{"valid package name", PackageName("main"), true},
		{"empty string", PackageName(""), false},
		{"whitespace only", PackageName("   "), false},
		{"package with slash", PackageName("package/name"), false},
		{"package with spaces", PackageName("package name"), true},
		{"single character", PackageName("p"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestPackageName_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     PackageName
		expectErr bool
	}{
		{"valid package name", PackageName("main"), false},
		{"empty string", PackageName(""), true},
		{"whitespace only", PackageName("   "), true},
		{"package with slash", PackageName("package/name"), true},
		{"package with spaces", PackageName("package name"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test custom error types
func TestUnsupportedReaderError(t *testing.T) {
	err := UnsupportedReaderError{ReaderType: "unknown"}
	expected := "unsupported reader type: unknown"
	assert.Equal(t, expected, err.Error())
}

func TestDirectoryCreationError(t *testing.T) {
	innerErr := errors.New("permission denied")
	err := DirectoryCreationError{Path: "/tmp/test", Err: innerErr}
	expected := "failed to create directory /tmp/test: permission denied"
	assert.Equal(t, expected, err.Error())
}

func TestFileWriteError(t *testing.T) {
	innerErr := errors.New("disk full")
	err := FileWriteError{Path: "/tmp/file.txt", Err: innerErr}
	expected := "failed to write /tmp/file.txt: disk full"
	assert.Equal(t, expected, err.Error())
}

func TestFileFormatError(t *testing.T) {
	innerErr := errors.New("invalid syntax")
	err := FileFormatError{Path: "/tmp/code.go", Err: innerErr}
	expected := "failed to format /tmp/code.go: invalid syntax"
	assert.Equal(t, expected, err.Error())
}

func TestModuleNotFoundError(t *testing.T) {
	err := ModuleNotFoundError{Source: "go.mod"}
	expected := "module declaration not found in go.mod"
	assert.Equal(t, expected, err.Error())
}

// Test Package struct
func TestPackage(t *testing.T) {
	pkg := Package{
		Path:  LayerPath("domain/entity"),
		Level: 1,
	}

	assert.Equal(t, "domain/entity", pkg.Path.String())
	assert.Equal(t, 1, pkg.Level)
}

// Test Layer struct
func TestLayer(t *testing.T) {
	layer := Layer{
		Name:  LayerName("Domain layer"),
		Order: 1,
		Packages: []Package{
			{Path: LayerPath("domain/entity"), Level: 0},
			{Path: LayerPath("domain/service"), Level: 1},
		},
	}

	assert.Equal(t, "Domain layer", layer.Name.String())
	assert.Equal(t, 1, layer.Order)
	assert.Len(t, layer.Packages, 2)
	assert.Equal(t, "domain/entity", layer.Packages[0].Path.String())
	assert.Equal(t, 0, layer.Packages[0].Level)
	assert.Equal(t, "domain/service", layer.Packages[1].Path.String())
	assert.Equal(t, 1, layer.Packages[1].Level)
}

// Test DependencyConfig struct
func TestDependencyConfig(t *testing.T) {
	config := DependencyConfig{
		Layers: []Layer{
			{
				Name:  LayerName("Domain layer"),
				Order: 1,
				Packages: []Package{
					{Path: LayerPath("domain/entity"), Level: 0},
				},
			},
			{
				Name:  LayerName("Application layer"),
				Order: 2,
				Packages: []Package{
					{Path: LayerPath("app/usecase"), Level: 0},
				},
			},
		},
	}

	assert.Len(t, config.Layers, 2)
	assert.Equal(t, "Domain layer", config.Layers[0].Name.String())
	assert.Equal(t, "Application layer", config.Layers[1].Name.String())
}

// Test GetAllPackages method
func TestDependencyConfig_GetAllPackages(t *testing.T) {
	config := &DependencyConfig{
		Layers: []Layer{
			{
				Name:  LayerName("Domain layer"),
				Order: 1,
				Packages: []Package{
					{Path: LayerPath("domain/entity"), Level: 0},
					{Path: LayerPath("domain/service"), Level: 1},
				},
			},
			{
				Name:  LayerName("Application layer"),
				Order: 2,
				Packages: []Package{
					{Path: LayerPath("app/usecase"), Level: 0},
				},
			},
		},
	}

	packages := config.GetAllPackages()

	assert.Len(t, packages, 3)

	expectedPaths := []string{"domain/entity", "domain/service", "app/usecase"}
	actualPaths := make([]string, len(packages))
	for i, pkg := range packages {
		actualPaths[i] = pkg.Path.String()
	}

	for _, expectedPath := range expectedPaths {
		assert.Contains(t, actualPaths, expectedPath)
	}
}

// Test GetPackagesByLayer method
func TestDependencyConfig_GetPackagesByLayer(t *testing.T) {
	config := &DependencyConfig{
		Layers: []Layer{
			{
				Name:  LayerName("Domain layer"),
				Order: 1,
				Packages: []Package{
					{Path: LayerPath("domain/entity"), Level: 0},
					{Path: LayerPath("domain/service"), Level: 1},
				},
			},
			{
				Name:  LayerName("Application layer"),
				Order: 2,
				Packages: []Package{
					{Path: LayerPath("app/usecase"), Level: 0},
				},
			},
		},
	}

	// Test existing layer
	domainPackages := config.GetPackagesByLayer(LayerName("Domain layer"))
	assert.Len(t, domainPackages, 2)
	assert.Equal(t, "domain/entity", domainPackages[0].Path.String())
	assert.Equal(t, "domain/service", domainPackages[1].Path.String())

	// Test non-existent layer
	nonExistentPackages := config.GetPackagesByLayer(LayerName("Non-existent layer"))
	assert.Nil(t, nonExistentPackages)
}

// Test GetDependenciesForPackage method
func TestDependencyConfig_GetDependenciesForPackage(t *testing.T) {
	config := &DependencyConfig{
		Layers: []Layer{
			{
				Name:  LayerName("Domain layer"),
				Order: 1,
				Packages: []Package{
					{Path: LayerPath("domain/entity"), Level: 0},
					{Path: LayerPath("domain/valueobject"), Level: 0},
					{Path: LayerPath("domain/service"), Level: 1},
				},
			},
			{
				Name:  LayerName("Application layer"),
				Order: 2,
				Packages: []Package{
					{Path: LayerPath("app/service"), Level: 0},
					{Path: LayerPath("app/usecase"), Level: 1},
				},
			},
		},
	}

	tests := []struct {
		name          string
		targetPackage Package
		expectedDeps  []string
	}{
		{
			name:          "domain entity - no dependencies",
			targetPackage: Package{Path: LayerPath("domain/entity"), Level: 0},
			expectedDeps:  []string{},
		},
		{
			name:          "domain service - depends on entity and valueobject",
			targetPackage: Package{Path: LayerPath("domain/service"), Level: 1},
			expectedDeps:  []string{"domain/entity", "domain/valueobject"},
		},
		{
			name:          "app usecase - depends on domain layer and app service",
			targetPackage: Package{Path: LayerPath("app/usecase"), Level: 1},
			expectedDeps:  []string{"domain/entity", "domain/valueobject", "domain/service", "app/service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dependencies := config.GetDependenciesForPackage(tt.targetPackage)

			actualDeps := make([]string, len(dependencies))
			for i, dep := range dependencies {
				actualDeps[i] = dep.String()
			}

			assert.Len(t, actualDeps, len(tt.expectedDeps))
			for _, expectedDep := range tt.expectedDeps {
				assert.Contains(t, actualDeps, expectedDep)
			}
		})
	}
}

// Test GetPackageName function
func TestGetPackageName(t *testing.T) {
	tests := []struct {
		name        string
		packagePath LayerPath
		expected    PackageName
	}{
		{
			name:        "simple package name",
			packagePath: LayerPath("entity"),
			expected:    PackageName("entity"),
		},
		{
			name:        "nested package path",
			packagePath: LayerPath("domain/entity"),
			expected:    PackageName("entity"),
		},
		{
			name:        "deeply nested package path",
			packagePath: LayerPath("infra/database/mysql"),
			expected:    PackageName("mysql"),
		},
		{
			name:        "package with multiple segments",
			packagePath: LayerPath("app/service/user"),
			expected:    PackageName("user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPackageName(tt.packagePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}
