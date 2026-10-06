package inkssg

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

//go:embed themes/*/*.html themes/*/*.css themes/*/*.js
var embeddedThemes embed.FS

var themeFS embed.FS = embeddedThemes
var hasThemes bool = true

func hasEmbeddedThemes() bool {
	return hasThemes
}

type Site struct {
	Dir      string
	Pages    []Page
	Output   string
	Theme    string
	PagesDir string
	Config   *SiteConfig
}

type Page struct {
	Name        string
	Dir         string
	Content     string
	Title       string
	Description string
	Theme       string
	ContentHTML string
}

type SiteConfig struct {
	Name         string    `yaml:"name"`
	Avatar       string    `yaml:"avatar"`
	Bio          string    `yaml:"bio"`
	Status       string    `yaml:"status"`
	Install      string    `yaml:"install"`
	DefaultTheme string    `yaml:"default_theme"`
	OutputDir    string    `yaml:"output_dir"`
	PagesDir     string    `yaml:"pages_dir"`
	Links        []Link    `yaml:"links"`
	Projects     []Project `yaml:"projects"`
	Meta         Meta      `yaml:"meta"`
	Contact      Contact   `yaml:"contact"`
}

type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
	Image       string `yaml:"image"`
	Badge       string `yaml:"badge"`
}

type Link struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Meta struct {
	Lang        string `yaml:"lang"`
	Description string `yaml:"description"`
	Title       string `yaml:"title"`
	Author      string `yaml:"author"`
	SiteUrl     string `yaml:"siteUrl"`
}

type Contact struct {
	Socials []Link `yaml:"socials"`
}

type TemplateData struct {
	Site    *SiteConfig
	Page    Page
	Theme   map[string]string
	Content template.HTML
}

func Scaffold(dir string) error {
	if dir == "" {
		dir = "."
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	pagesDir := filepath.Join(absDir, "pages")
	if _, err := os.Stat(pagesDir); err == nil {
		return fmt.Errorf("pages/ already exists at %s — refusing to overwrite", absDir)
	}

	indexDir := filepath.Join(pagesDir, "index")
	if err := os.MkdirAll(indexDir, 0755); err != nil {
		return fmt.Errorf("creating pages/index: %w", err)
	}

	indexContent := `---
title: Hello
description: My new inkssg site
---

# Hello

Edit ` + "`pages/index/content.md`" + ` and run ` + "`inkssg build`" + `.
`
	if err := os.WriteFile(filepath.Join(indexDir, "content.md"), []byte(indexContent), 0644); err != nil {
		return fmt.Errorf("writing content.md: %w", err)
	}

	configPath := filepath.Join(absDir, "ink.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		sample := `name: My Site
default_theme: minimal
output_dir: public

meta:
  lang: en
  title: My Site
  description: A small site built with inkssg
`
		if err := os.WriteFile(configPath, []byte(sample), 0644); err != nil {
			return fmt.Errorf("writing ink.yaml: %w", err)
		}
	}

	fmt.Printf("scaffolded site at %s\n", absDir)
	fmt.Println("next: inkssg build")
	return nil
}

func Build(paths ...string) error {
	dir := "."
	if len(paths) > 0 {
		dir = paths[0]
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	site, err := NewSite(absDir)
	if err != nil {
		return err
	}

	return site.Build()
}

func NewSite(dir string) (*Site, error) {
	site := &Site{
		Dir:      dir,
		Pages:    []Page{},
		Output:   "public",
		Theme:    "minimal",
		PagesDir: "pages",
		Config:   &SiteConfig{},
	}

	if err := site.loadConfig(); err != nil {
		// config is optional
	}

	if err := site.detect(); err != nil {
		return nil, err
	}

	return site, nil
}

func (s *Site) loadConfig() error {
	configPath := filepath.Join(s.Dir, "ink.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	var config SiteConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("invalid ink.yaml: %w", err)
	}

	if config.OutputDir != "" {
		s.Output = config.OutputDir
	}

	if config.DefaultTheme != "" {
		s.Theme = config.DefaultTheme
	}

	if config.PagesDir != "" {
		s.PagesDir = config.PagesDir
	}

	s.Config = &config

	return nil
}

func (s *Site) detect() error {
	pagesDir := filepath.Join(s.Dir, s.PagesDir)
	if _, err := os.Stat(pagesDir); os.IsNotExist(err) {
		return fmt.Errorf("no %s/ folder found. Run 'inkssg new .' to scaffold a site.", s.PagesDir)
	}

	entries, err := os.ReadDir(pagesDir)
	if err != nil {
		return fmt.Errorf("cannot read %s/: %w", s.PagesDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pageDir := filepath.Join(pagesDir, entry.Name())
		contentPath := filepath.Join(pageDir, "content.md")
		contentType := "md"
		if _, err := os.Stat(contentPath); os.IsNotExist(err) {
			contentPath = filepath.Join(pageDir, "content.html")
			contentType = "html"
			if _, err := os.Stat(contentPath); os.IsNotExist(err) {
				continue
			}
		}

		page := Page{
			Name:    entry.Name(),
			Dir:     pageDir,
			Content: contentPath,
		}

		if err := page.parseFrontmatter(contentType); err != nil {
			return fmt.Errorf("%s: %w", entry.Name(), err)
		}

		if page.Theme == "" {
			page.Theme = s.Theme
		}

		s.Pages = append(s.Pages, page)
	}

	if len(s.Pages) == 0 {
		return fmt.Errorf("no pages found. Add a pages/<name>/content.md or content.html file.")
	}

	return nil
}

func (p *Page) parseFrontmatter(contentType string) error {
	data, err := os.ReadFile(p.Content)
	if err != nil {
		return fmt.Errorf("cannot read content: %w", err)
	}

	frontmatter := &struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Theme       string `yaml:"theme"`
	}{}

	str := string(data)

	if strings.HasPrefix(str, "---") {
		str = strings.TrimPrefix(str, "---\n")
		endIdx := strings.Index(str, "\n---")
		if endIdx > 0 {
			yamlData := str[:endIdx]
			if err := yaml.Unmarshal([]byte(yamlData), frontmatter); err != nil {
				return fmt.Errorf("invalid frontmatter: %w", err)
			}
			content := str[endIdx+4:]
			p.ContentHTML = p.renderMarkdown([]byte(content), contentType)
		} else {
			p.ContentHTML = p.renderMarkdown(data, contentType)
		}
	} else {
		p.ContentHTML = p.renderMarkdown(data, contentType)
	}

	p.Title = frontmatter.Title
	p.Description = frontmatter.Description
	p.Theme = frontmatter.Theme

	if p.Title == "" {
		p.Title = p.Name
	}

	return nil
}

func (p *Page) renderMarkdown(data []byte, contentType string) string {
	if contentType == "html" {
		return string(data)
	}

	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert(data, &buf); err != nil {
		return ""
	}
	return buf.String()
}

func (s *Site) validateOutput() error {
	if strings.TrimSpace(s.Output) == "" {
		return fmt.Errorf("output_dir cannot be empty")
	}

	absDir, err := filepath.Abs(s.Dir)
	if err != nil {
		return fmt.Errorf("invalid project dir: %w", err)
	}
	absOut, err := filepath.Abs(filepath.Join(s.Dir, s.Output))
	if err != nil {
		return fmt.Errorf("invalid output_dir: %w", err)
	}

	if absOut == absDir {
		return fmt.Errorf("output_dir cannot be the project root (would delete sources on rebuild)")
	}

	rel, err := filepath.Rel(absDir, absOut)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("output_dir must be inside the project directory, got %q", s.Output)
	}

	return nil
}

func (s *Site) Build() error {
	start := time.Now()

	if err := s.validateOutput(); err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Join(s.Dir, s.Output)); err != nil {
		return fmt.Errorf("cannot clean output dir: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(s.Dir, s.Output), 0755); err != nil {
		return fmt.Errorf("cannot create output dir: %w", err)
	}

	themeDir, err := s.resolveTheme()
	if err != nil {
		return err
	}

	if err := s.copyTheme(themeDir); err != nil {
		return fmt.Errorf("copying theme: %w", err)
	}

	if err := s.copyAssets(); err != nil {
		return fmt.Errorf("copying assets: %w", err)
	}

	if err := s.copyStatic(); err != nil {
		return fmt.Errorf("copying static: %w", err)
	}

	failed := 0
	for _, page := range s.Pages {
		if err := s.buildPage(page, themeDir); err != nil {
			fmt.Fprintf(os.Stderr, "✗ %s: %v\n", page.Name, err)
			failed++
		} else {
			fmt.Printf("✓ %s → %s.html\n", page.Name, page.Name)
		}
	}

	total := len(s.Pages)
	duration := time.Since(start)

	if failed > 0 {
		fmt.Printf("%d page failed (of %d) in %s\n", failed, total, duration)
		return fmt.Errorf("build failed")
	}

	fmt.Printf("built %d pages in %s\n", total, duration)
	return nil
}

func (s *Site) resolveTheme() (string, error) {
	localTheme := filepath.Join(s.Dir, "themes", s.Theme)
	if _, err := os.Stat(localTheme); err == nil {
		return localTheme, nil
	}

	if hasEmbeddedThemes() {
		return "", nil // use embedded
	}

	return "", fmt.Errorf("theme %q not found (local or built-in)", s.Theme)
}

func (s *Site) copyTheme(themeDir string) error {
	outputTheme := filepath.Join(s.Dir, s.Output, "themes", s.Theme)
	if err := os.MkdirAll(outputTheme, 0755); err != nil {
		return err
	}

	if themeDir == "" && hasEmbeddedThemes() {
		files := []string{"styles.css", "script.js"}
		for _, f := range files {
			data, err := themeFS.ReadFile("themes/" + s.Theme + "/" + f)
			if err == nil {
				os.WriteFile(filepath.Join(outputTheme, f), data, 0644)
			}
		}
		return nil
	}

	entries, _ := os.ReadDir(themeDir)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == "layout.html" {
			continue
		}
		src := filepath.Join(themeDir, entry.Name())
		dst := filepath.Join(outputTheme, entry.Name())
		data, _ := os.ReadFile(src)
		os.WriteFile(dst, data, 0644)
	}

	return nil
}

func (s *Site) copyAssets() error {
	assetsDir := filepath.Join(s.Dir, "assets")
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		return nil
	}

	outputAssets := filepath.Join(s.Dir, s.Output, "assets")
	if err := os.MkdirAll(outputAssets, 0755); err != nil {
		return err
	}

	return copyDir(assetsDir, outputAssets)
}

// copyStatic copies static/ to the root of the output, so static/cv.pdf is
// served at /cv.pdf. It runs before pages are written: a page wins over a
// static file with the same name.
func (s *Site) copyStatic() error {
	staticDir := filepath.Join(s.Dir, "static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		return nil
	}

	return copyDir(staticDir, filepath.Join(s.Dir, s.Output))
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			os.MkdirAll(dstPath, 0755)
			copyDir(srcPath, dstPath)
		} else {
			data, _ := os.ReadFile(srcPath)
			os.WriteFile(dstPath, data, 0644)
		}
	}

	return nil
}

func (s *Site) buildPage(page Page, themeDir string) error {
	var layout string
	themeName := page.Theme
	if themeName == "" {
		themeName = s.Theme
	}

	if themeDir == "" && hasEmbeddedThemes() {
		data, err := themeFS.ReadFile("themes/" + themeName + "/layout.html")
		if err != nil {
			return fmt.Errorf("cannot read embedded layout for theme %q: %w", themeName, err)
		}
		layout = string(data)
	} else if themeDir != "" {
		data, err := os.ReadFile(filepath.Join(themeDir, "layout.html"))
		if err != nil {
			return fmt.Errorf("cannot read layout: %w", err)
		}
		layout = string(data)
	} else {
		return fmt.Errorf("no layout found for theme %q", themeName)
	}

	data := TemplateData{
		Site:    s.Config,
		Page:    page,
		Theme:   nil,
		Content: template.HTML(page.ContentHTML),
	}

	tmpl, err := template.New("layout").Parse(layout)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	outputPath := filepath.Join(s.Dir, s.Output, page.Name+".html")
	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
