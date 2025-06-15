package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDependencyContent(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectedLayers []Layer
		expectError    bool
	}{
		{
			name: "valid hierarchical dependency file",
			content: `# Dependencies

## Layers

Upper layers cannot depend on lower layers.

1. Domain layer
  - Implementation of core entities
2. Application layer
  - Business logic using objects from the domain layer
3. Presentation layer
  - UI and presentation logic
4. Infra layer
  - Gateway to the real world

## Packages in layers

Upper packages cannot depend on lower packages.

1. Domain layer
  - domain/entity
  - domain/valueobject
    - domain/service
2. Application layer
  - app/service
    - app/usecase
3. Presentation layer
  - api
  - cli
4. Infra layer
  - infra/database
  - infra/cache
`,
			expectedLayers: []Layer{
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
				{
					Name:  LayerName("Presentation layer"),
					Order: 3,
					Packages: []Package{
						{Path: LayerPath("api"), Level: 0},
						{Path: LayerPath("cli"), Level: 0},
					},
				},
				{
					Name:  LayerName("Infra layer"),
					Order: 4,
					Packages: []Package{
						{Path: LayerPath("infra/database"), Level: 0},
						{Path: LayerPath("infra/cache"), Level: 0},
					},
				},
			},
			expectError: false,
		},
		{
			name:           "empty content",
			content:        ``,
			expectedLayers: []Layer{},
			expectError:    false,
		},
		{
			name: "only layers section",
			content: `## Layers
1. Test layer
  - Test layer description`,
			expectedLayers: []Layer{
				{Name: LayerName("Test layer"), Order: 1, Packages: []Package{}},
			},
			expectError: false,
		},
		{
			name: "only packages section",
			content: `## Packages in layers
1. Domain layer
  - domain/entity`,
			expectedLayers: []Layer{},
			expectError:    false,
		},
		{
			name: "invalid layer path with ..",
			content: `## Layers
1. Test layer

## Packages in layers
1. Test layer
  - ../invalid/path`,
			expectedLayers: []Layer{},
			expectError:    true,
		},
		{
			name: "single layer with single package",
			content: `## Layers
1. Domain layer

## Packages in layers
1. Domain layer
  - domain/entity`,
			expectedLayers: []Layer{
				{
					Name:     LayerName("Domain layer"),
					Order:    1,
					Packages: []Package{{Path: LayerPath("domain/entity"), Level: 0}},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			config, err := parser.ParseDependencyContent(tt.content)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, config)

			// Check layers
			assert.Len(t, config.Layers, len(tt.expectedLayers))

			for i, expected := range tt.expectedLayers {
				require.Less(t, i, len(config.Layers), "Missing layer at index %d", i)
				actual := config.Layers[i]
				assert.Equal(t, expected.Name, actual.Name, "Layer %d name mismatch", i)
				assert.Equal(t, expected.Order, actual.Order, "Layer %d order mismatch", i)

				// Check packages
				assert.Len(t, actual.Packages, len(expected.Packages), "Layer %d package count mismatch", i)
				for j, expectedPkg := range expected.Packages {
					require.Less(t, j, len(actual.Packages), "Missing package at index %d for layer %d", j, i)
					actualPkg := actual.Packages[j]
					assert.Equal(t, expectedPkg.Path, actualPkg.Path, "Layer %d package %d path mismatch", i, j)
					assert.Equal(t, expectedPkg.Level, actualPkg.Level, "Layer %d package %d level mismatch", i, j)
				}
			}
		})
	}
}

func TestParseLayersSection(t *testing.T) {
	tests := []struct {
		name           string
		lines          []string
		expectedLayers []Layer
		expectError    bool
	}{
		{
			name: "single layer",
			lines: []string{
				"1. Domain layer",
			},
			expectedLayers: []Layer{
				{Name: LayerName("Domain layer"), Order: 1, Packages: []Package{}},
			},
			expectError: false,
		},
		{
			name: "multiple layers",
			lines: []string{
				"1. Domain layer",
				"2. Application layer",
				"3. Presentation layer",
			},
			expectedLayers: []Layer{
				{Name: LayerName("Domain layer"), Order: 1, Packages: []Package{}},
				{Name: LayerName("Application layer"), Order: 2, Packages: []Package{}},
				{Name: LayerName("Presentation layer"), Order: 3, Packages: []Package{}},
			},
			expectError: false,
		},
		{
			name: "layers with descriptions",
			lines: []string{
				"1. Domain layer",
				"  - Implementation of core entities",
				"2. Application layer",
				"  - Business logic",
			},
			expectedLayers: []Layer{
				{Name: LayerName("Domain layer"), Order: 1, Packages: []Package{}},
				{Name: LayerName("Application layer"), Order: 2, Packages: []Package{}},
			},
			expectError: false,
		},
		{
			name: "invalid layer number",
			lines: []string{
				"x. Invalid layer",
			},
			expectedLayers: []Layer{},
			expectError:    false, // Should skip invalid lines
		},
		{
			name: "empty layer name",
			lines: []string{
				"1. ",
			},
			expectedLayers: []Layer{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DependencyConfig{
				Layers: make([]Layer, 0),
			}

			parser := NewParser()
			var lastErr error
			for _, line := range tt.lines {
				err := parser.ParseLayersSection(line, config)
				if err != nil {
					lastErr = err
					break
				}
			}

			if tt.expectError {
				assert.Error(t, lastErr)
				return
			}

			assert.NoError(t, lastErr)
			assert.Len(t, config.Layers, len(tt.expectedLayers))

			for i, expected := range tt.expectedLayers {
				require.Less(t, i, len(config.Layers), "Missing layer at index %d", i)
				actual := config.Layers[i]
				assert.Equal(t, expected.Name, actual.Name, "Layer %d name mismatch", i)
				assert.Equal(t, expected.Order, actual.Order, "Layer %d order mismatch", i)
			}
		})
	}
}

func TestParsePackagesSection(t *testing.T) {
	tests := []struct {
		name            string
		setupLayers     []Layer
		lines           []string
		expectedPackage map[string][]Package // layer name -> packages
		expectError     bool
	}{
		{
			name: "single layer with packages",
			setupLayers: []Layer{
				{Name: LayerName("Domain layer"), Order: 1, Packages: []Package{}},
			},
			lines: []string{
				"1. Domain layer",
				"  - domain/entity",
				"  - domain/valueobject",
				"    - domain/service",
			},
			expectedPackage: map[string][]Package{
				"Domain layer": {
					{Path: LayerPath("domain/entity"), Level: 0},
					{Path: LayerPath("domain/valueobject"), Level: 0},
					{Path: LayerPath("domain/service"), Level: 1},
				},
			},
			expectError: false,
		},
		{
			name: "multiple layers with packages",
			setupLayers: []Layer{
				{Name: LayerName("Domain layer"), Order: 1, Packages: []Package{}},
				{Name: LayerName("Application layer"), Order: 2, Packages: []Package{}},
			},
			lines: []string{
				"1. Domain layer",
				"  - domain/entity",
				"2. Application layer",
				"  - app/service",
				"    - app/usecase",
			},
			expectedPackage: map[string][]Package{
				"Domain layer": {
					{Path: LayerPath("domain/entity"), Level: 0},
				},
				"Application layer": {
					{Path: LayerPath("app/service"), Level: 0},
					{Path: LayerPath("app/usecase"), Level: 1},
				},
			},
			expectError: false,
		},
		{
			name: "invalid package path",
			setupLayers: []Layer{
				{Name: LayerName("Domain layer"), Order: 1, Packages: []Package{}},
			},
			lines: []string{
				"1. Domain layer",
				"  - ../invalid/path",
			},
			expectedPackage: map[string][]Package{},
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DependencyConfig{
				Layers: tt.setupLayers,
			}

			parser := NewParser()
			var currentLayer *Layer
			var lastErr error

			for _, line := range tt.lines {
				err := parser.ParsePackagesSection(line, config, &currentLayer)
				if err != nil {
					lastErr = err
					break
				}
			}

			if tt.expectError {
				assert.Error(t, lastErr)
				return
			}

			assert.NoError(t, lastErr)

			// Check packages in each layer
			for layerName, expectedPackages := range tt.expectedPackage {
				var foundLayer *Layer
				for i := range config.Layers {
					if string(config.Layers[i].Name) == layerName {
						foundLayer = &config.Layers[i]
						break
					}
				}

				require.NotNil(t, foundLayer, "Layer %s not found", layerName)
				assert.Len(t, foundLayer.Packages, len(expectedPackages), "Package count mismatch for layer %s", layerName)

				for i, expectedPkg := range expectedPackages {
					require.Less(t, i, len(foundLayer.Packages), "Missing package at index %d for layer %s", i, layerName)
					actualPkg := foundLayer.Packages[i]
					assert.Equal(t, expectedPkg.Path, actualPkg.Path, "Package %d path mismatch for layer %s", i, layerName)
					assert.Equal(t, expectedPkg.Level, actualPkg.Level, "Package %d level mismatch for layer %s", i, layerName)
				}
			}
		})
	}
}

func TestCalculateIndentationLevel(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected int
	}{
		{
			name:     "no indentation",
			line:     "- package/path",
			expected: 0,
		},
		{
			name:     "level 0 indentation",
			line:     "  - package/path",
			expected: 0,
		},
		{
			name:     "level 1 indentation",
			line:     "    - package/path",
			expected: 1,
		},
		{
			name:     "extra spaces",
			line:     "      - package/path",
			expected: 1,
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.calculateIndentationLevel(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetModuleNameFromContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    ModuleName
		expectError bool
	}{
		{
			name: "valid module declaration",
			content: `module github.com/example/project

go 1.21
`,
			expected:    ModuleName("github.com/example/project"),
			expectError: false,
		},
		{
			name: "module with extra spaces",
			content: `module   github.com/example/project

go 1.21
`,
			expected:    ModuleName("github.com/example/project"),
			expectError: false,
		},
		{
			name: "module declaration in middle of file",
			content: `// This is a go.mod file
module github.com/test/module

require (
    github.com/dependency v1.0.0
)
`,
			expected:    ModuleName("github.com/test/module"),
			expectError: false,
		},
		{
			name: "no module declaration",
			content: `go 1.21

require (
    github.com/dependency v1.0.0
)
`,
			expected:    ModuleName(""),
			expectError: true,
		},
		{
			name:        "empty content",
			content:     ``,
			expected:    ModuleName(""),
			expectError: true,
		},
		{
			name: "invalid module name with spaces",
			content: `module invalid module name

go 1.21
`,
			expected:    ModuleName(""),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			result, err := parser.GetModuleNameFromContent(tt.content, "test.mod")

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseDependencyFile(t *testing.T) {
	// Create a temporary file
	content := `# Dependencies

## Layers

Upper layers cannot depend on lower layers.

1. Domain layer
  - Implementation of core entities
2. Application layer
  - Business logic using objects from the domain layer

## Packages in layers

Upper packages cannot depend on lower packages.

1. Domain layer
  - domain/entity
  - domain/valueobject
2. Application layer
  - app/service
    - app/usecase
`

	tmpFile, err := os.CreateTemp("", "dependency-*.md")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	tmpFile.Close()

	parser := NewParser()
	config, err := parser.ParseDependencyFile(tmpFile.Name())
	require.NoError(t, err)

	// Verify the parsed content
	expectedLayers := []Layer{
		{
			Name:  LayerName("Domain layer"),
			Order: 1,
			Packages: []Package{
				{Path: LayerPath("domain/entity"), Level: 0},
				{Path: LayerPath("domain/valueobject"), Level: 0},
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
	}

	assert.Len(t, config.Layers, len(expectedLayers))

	for i, expected := range expectedLayers {
		require.Less(t, i, len(config.Layers), "Missing layer at index %d", i)
		actual := config.Layers[i]
		assert.Equal(t, expected.Name, actual.Name, "Layer %d name mismatch", i)
		assert.Equal(t, expected.Order, actual.Order, "Layer %d order mismatch", i)

		assert.Len(t, actual.Packages, len(expected.Packages), "Layer %d package count mismatch", i)
		for j, expectedPkg := range expected.Packages {
			require.Less(t, j, len(actual.Packages), "Missing package at index %d for layer %d", j, i)
			actualPkg := actual.Packages[j]
			assert.Equal(t, expectedPkg.Path, actualPkg.Path, "Layer %d package %d path mismatch", i, j)
			assert.Equal(t, expectedPkg.Level, actualPkg.Level, "Layer %d package %d level mismatch", i, j)
		}
	}
}

func TestGetModuleName(t *testing.T) {
	// Create a temporary go.mod file
	content := `module github.com/test/project

go 1.21

require (
    github.com/dependency v1.0.0
)
`

	tmpFile, err := os.CreateTemp("", "go.mod")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	tmpFile.Close()

	parser := NewParser()
	result, err := parser.GetModuleName(tmpFile.Name())
	require.NoError(t, err)

	expected := ModuleName("github.com/test/project")
	assert.Equal(t, expected, result)
}
