package main

import (
	"fmt"
	"strings"
)

// Custom types for better type safety
type LayerName string
type LayerPath string
type ModuleName string
type PackageName string
type FilePath string

func (ln LayerName) String() string   { return string(ln) }
func (lp LayerPath) String() string   { return string(lp) }
func (mn ModuleName) String() string  { return string(mn) }
func (pn PackageName) String() string { return string(pn) }
func (fp FilePath) String() string    { return string(fp) }

// Validation methods for custom types
func (ln LayerName) IsValid() bool {
	return ln != "" && strings.TrimSpace(string(ln)) != ""
}

func (ln LayerName) Validate() error {
	if !ln.IsValid() {
		return fmt.Errorf("layer name cannot be empty")
	}
	return nil
}

func (lp LayerPath) IsValid() bool {
	return lp != "" && strings.TrimSpace(string(lp)) != "" && !strings.Contains(string(lp), "..")
}

func (lp LayerPath) Validate() error {
	if lp == "" || strings.TrimSpace(string(lp)) == "" {
		return fmt.Errorf("layer path cannot be empty")
	}
	if strings.Contains(string(lp), "..") {
		return fmt.Errorf("layer path cannot contain '..' for security reasons")
	}
	return nil
}

func (mn ModuleName) IsValid() bool {
	return mn != "" && strings.TrimSpace(string(mn)) != "" && !strings.Contains(string(mn), " ")
}

func (mn ModuleName) Validate() error {
	if mn == "" || strings.TrimSpace(string(mn)) == "" {
		return fmt.Errorf("module name cannot be empty")
	}
	if strings.Contains(string(mn), " ") {
		return fmt.Errorf("module name cannot contain spaces")
	}
	return nil
}

func (pn PackageName) IsValid() bool {
	return pn != "" && strings.TrimSpace(string(pn)) != "" && !strings.Contains(string(pn), "/")
}

func (pn PackageName) Validate() error {
	if pn == "" || strings.TrimSpace(string(pn)) == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	if strings.Contains(string(pn), "/") {
		return fmt.Errorf("package name cannot contain '/'")
	}
	return nil
}

// Custom error types for better error handling
type UnsupportedReaderError struct {
	ReaderType string
}

func (e UnsupportedReaderError) Error() string {
	return fmt.Sprintf("unsupported reader type: %s", e.ReaderType)
}

type DirectoryCreationError struct {
	Path string
	Err  error
}

func (e DirectoryCreationError) Error() string {
	return fmt.Sprintf("failed to create directory %s: %v", e.Path, e.Err)
}

type FileWriteError struct {
	Path string
	Err  error
}

func (e FileWriteError) Error() string {
	return fmt.Sprintf("failed to write %s: %v", e.Path, e.Err)
}

type FileFormatError struct {
	Path string
	Err  error
}

func (e FileFormatError) Error() string {
	return fmt.Sprintf("failed to format %s: %v", e.Path, e.Err)
}

type ModuleNotFoundError struct {
	Source string
}

func (e ModuleNotFoundError) Error() string {
	return fmt.Sprintf("module declaration not found in %s", e.Source)
}

// Package represents a single package with its hierarchical level
type Package struct {
	Path  LayerPath // e.g., "domain/entity", "domain/service"
	Level int       // Indentation level (0 = top level, 1 = one indent, etc.)
}

// Layer represents a layer with its packages
type Layer struct {
	Name     LayerName
	Order    int       // Layer order (1, 2, 3, ...)
	Packages []Package // Packages in this layer
}

// DependencyConfig represents the complete dependency configuration
type DependencyConfig struct {
	Layers []Layer
}

// GetAllPackages returns all packages across all layers
func (dc *DependencyConfig) GetAllPackages() []Package {
	var allPackages []Package
	for _, layer := range dc.Layers {
		allPackages = append(allPackages, layer.Packages...)
	}
	return allPackages
}

// GetPackagesByLayer returns packages for a specific layer
func (dc *DependencyConfig) GetPackagesByLayer(layerName LayerName) []Package {
	for _, layer := range dc.Layers {
		if layer.Name == layerName {
			return layer.Packages
		}
	}
	return nil
}

// GetDependenciesForPackage calculates dependencies for a given package
func (dc *DependencyConfig) GetDependenciesForPackage(targetPackage Package) []LayerPath {
	var dependencies []LayerPath

	// Find the layer containing this package
	var targetLayer *Layer
	for i := range dc.Layers {
		for _, pkg := range dc.Layers[i].Packages {
			if pkg.Path == targetPackage.Path {
				targetLayer = &dc.Layers[i]
				break
			}
		}
		if targetLayer != nil {
			break
		}
	}

	if targetLayer == nil {
		return dependencies
	}

	// Add dependencies from upper layers (layers with lower order)
	for _, layer := range dc.Layers {
		if layer.Order < targetLayer.Order {
			for _, pkg := range layer.Packages {
				dependencies = append(dependencies, pkg.Path)
			}
		}
	}

	// Find target package index
	var targetIndex int
	for i, pkg := range targetLayer.Packages {
		if pkg.Path == targetPackage.Path {
			targetIndex = i
			break
		}
	}

	// Add dependencies from the same layer
	for i, pkg := range targetLayer.Packages {
		if pkg.Path != targetPackage.Path {
			// Can depend on packages at higher levels (lower level number)
			// or packages at the same level that come before in the hierarchy
			if pkg.Level < targetPackage.Level || (pkg.Level == targetPackage.Level && i < targetIndex) {
				dependencies = append(dependencies, pkg.Path)
			}
		}
	}

	return dependencies
}

// GetPackageName extracts the package name from a package path
func GetPackageName(packagePath LayerPath) PackageName {
	path := string(packagePath)
	parts := strings.Split(path, "/")
	return PackageName(parts[len(parts)-1])
}
