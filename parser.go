package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	emptyLayerRegex = regexp.MustCompile(`^(\d+)\.\s*$`)
	layerRegex      = regexp.MustCompile(`^(\d+)\.\s+(.+)$`)
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseDependencyFile(filePath string) (*DependencyConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return p.ParseDependencyContent(file)
}

func (p *Parser) ParseDependencyContent(reader any) (*DependencyConfig, error) {
	var scanner *bufio.Scanner

	switch r := reader.(type) {
	case *os.File:
		scanner = bufio.NewScanner(r)
	case *strings.Reader:
		scanner = bufio.NewScanner(r)
	case string:
		scanner = bufio.NewScanner(strings.NewReader(r))
	default:
		return nil, UnsupportedReaderError{ReaderType: fmt.Sprintf("%T", reader)}
	}

	config := &DependencyConfig{
		Layers: make([]Layer, 0),
	}

	inLayersSection := false
	inPackagesSection := false
	var currentLayer *Layer

	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.TrimSpace(rawLine)

		// Check for section headers
		if strings.HasPrefix(line, "## Layers") {
			inLayersSection = true
			inPackagesSection = false
			continue
		}
		if strings.HasPrefix(line, "## Packages in layers") {
			inLayersSection = false
			inPackagesSection = true
			currentLayer = nil
			continue
		}

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse layers section
		if inLayersSection {
			err := p.ParseLayersSection(rawLine, config)
			if err != nil {
				return nil, err
			}
		}

		// Parse packages section
		if inPackagesSection {
			err := p.ParsePackagesSection(rawLine, config, &currentLayer)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

func (p *Parser) ParseLayersSection(line string, config *DependencyConfig) error {
	trimmed := strings.TrimSpace(line)

	// Skip description lines and non-numbered lines
	if strings.HasPrefix(trimmed, "-") ||
		strings.Contains(trimmed, "cannot depend") ||
		trimmed == "" {
		return nil
	}

	// Check for numbered lines with empty names like "1. "
	// Using package-level emptyLayerRegex constant
	if emptyLayerRegex.MatchString(trimmed) {
		return fmt.Errorf("invalid layer name: layer name cannot be empty")
	}

	// Match numbered layer lines like "1. Domain layer"
	// Using package-level layerRegex constant
	matches := layerRegex.FindStringSubmatch(trimmed)

	if len(matches) == 3 {
		orderStr := matches[1]
		layerName := strings.TrimSpace(matches[2])

		order, err := strconv.Atoi(orderStr)
		if err != nil {
			return fmt.Errorf("invalid layer order: %s", orderStr)
		}

		layer := Layer{
			Name:     LayerName(layerName),
			Order:    order,
			Packages: make([]Package, 0),
		}

		if err := layer.Name.Validate(); err != nil {
			return fmt.Errorf("invalid layer name: %v", err)
		}

		config.Layers = append(config.Layers, layer)
	}

	return nil
}

func (p *Parser) ParsePackagesSection(line string, config *DependencyConfig, currentLayer **Layer) error {
	trimmed := strings.TrimSpace(line)

	// Skip description lines
	if strings.Contains(trimmed, "cannot depend") || trimmed == "" {
		return nil
	}

	// Check for numbered lines with empty names like "1. "
	// Using package-level emptyLayerRegex constant
	if emptyLayerRegex.MatchString(trimmed) {
		return fmt.Errorf("invalid layer name: layer name cannot be empty")
	}

	// Match numbered layer lines like "1. Domain layer"
	// Using package-level layerRegex constant
	matches := layerRegex.FindStringSubmatch(trimmed)

	if len(matches) == 3 {
		orderStr := matches[1]
		layerName := strings.TrimSpace(matches[2])

		order, err := strconv.Atoi(orderStr)
		if err != nil {
			return fmt.Errorf("invalid layer order: %s", orderStr)
		}

		// Find the corresponding layer in config
		for i := range config.Layers {
			if config.Layers[i].Name == LayerName(layerName) && config.Layers[i].Order == order {
				*currentLayer = &config.Layers[i]
				break
			}
		}

		return nil
	}

	// Match package lines with indentation
	if strings.HasPrefix(trimmed, "- ") && *currentLayer != nil {
		packagePath := strings.TrimPrefix(trimmed, "- ")
		packagePath = strings.TrimSpace(packagePath)

		if packagePath == "" {
			return nil
		}

		// Calculate indentation level
		level := p.calculateIndentationLevel(line)

		pkg := Package{
			Path:  LayerPath(packagePath),
			Level: level,
		}

		if err := pkg.Path.Validate(); err != nil {
			return fmt.Errorf("invalid package path: %v", err)
		}

		(*currentLayer).Packages = append((*currentLayer).Packages, pkg)
	}

	return nil
}

func (p *Parser) calculateIndentationLevel(line string) int {
	// Count leading spaces before the "- " marker
	spacesBeforeDash := 0
	for _, char := range line {
		if char == ' ' {
			spacesBeforeDash++
		} else if char == '-' {
			break
		} else {
			// Non-space, non-dash character before dash
			break
		}
	}

	// Convert spaces to indentation level using 4-space standard
	// This supports arbitrary nesting depth:
	// Level 0: 0-3 spaces
	// Level 1: 4-7 spaces ("    - ", "      - ")
	// Level 2: 8-11 spaces ("        - ")
	// Level n: 4*n to 4*n+3 spaces (supports unlimited depth)
	return spacesBeforeDash / 4
}

func (p *Parser) GetModuleName(goModPath string) (ModuleName, error) {
	file, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return p.GetModuleNameFromContent(file, goModPath)
}

func (p *Parser) GetModuleNameFromContent(reader any, sourceName string) (ModuleName, error) {
	var scanner *bufio.Scanner

	switch r := reader.(type) {
	case *os.File:
		scanner = bufio.NewScanner(r)
	case *strings.Reader:
		scanner = bufio.NewScanner(r)
	case string:
		scanner = bufio.NewScanner(strings.NewReader(r))
	default:
		return "", UnsupportedReaderError{ReaderType: fmt.Sprintf("%T", reader)}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			moduleName := ModuleName(strings.TrimSpace(strings.TrimPrefix(line, "module ")))
			if err := moduleName.Validate(); err != nil {
				return "", fmt.Errorf("invalid module name: %v", err)
			}
			return moduleName, nil
		}
	}

	return "", ModuleNotFoundError{Source: sourceName}
}
