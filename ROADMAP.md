# Roadmap

## v0.1 — Library-first DX

Goal: a user can `go get` inkssg and have a site running in minutes. No config needed for simple sites.

### Library API (primary DX)
- [x] `inkssg.Build()` — no args, defaults to current directory
- [x] `inkssg.Build(dir)` — explicit path (same behavior as CLI)
- [x] Auto-detect site structure (no `ink.yaml` required)
- [x] Defaults: `minimal` theme, `public/` output, `pages/` dir
- [x] Clear error if not a valid site

### CLI surface (wraps library)
- [x] `inkssg build` — calls `inkssg.Build(".")`
- [x] `inkssg new <path>` — scaffold from template
- [x] `inkssg version`

### Build core
- [x] Discover pages under `pages/<name>/`
- [x] Parse YAML frontmatter from `content.md` and `content.html`
- [x] Render markdown with goldmark
- [x] Pass HTML body through unchanged for `content.html`
- [x] Apply `themes/<name>/layout.html` via `html/template` with `{{.Site}}`, `{{.Page}}`, `{{.Content}}`
- [x] Output `<name>.html` at root of `public/`
- [x] Build all pages, report all errors, exit non-zero on any failure
- [x] Output: `✓ <page> → <file>` per page + summary line

### Themes
- [x] Theme resolution: local `themes/<name>/` → built-in (embedded) → error
- [x] Built-in `minimal` theme shipped via `go:embed`
- [x] Built-in `devtool` theme for dev tool landing pages
- [x] Theme requires only `layout.html`; styles/script optional
- [x] Copy theme `styles.css`, `script.js` to `public/themes/<name>/`

### Assets
- [x] Copy `assets/*` → `public/assets/` (one rule, no per-page folders)

### Examples + CI
- [x] `examples/library/` — `go run main.go` builds the site
- [x] `examples/minimal/` — 1 page, default theme
- [x] `examples/multi-page/` — multiple pages, shared theme
- [x] CI: build every example on push

Note: `ink.yaml` is optional in v0.1. Add it only when you need site-wide data (name, links, bio).

## v0.1.x — Fixes

Bugs and cleanup found in a code review. These come before new features.

### Bugs
- [ ] Reject `output_dir` set to `pages`, `assets`, `static` or `themes` (the build deletes them today)
- [ ] Fail the build on an invalid `ink.yaml` instead of falling back to defaults
- [ ] `devtool` theme links `/themes/landing/styles.css`, so it ships unstyled
- [ ] Per-page `theme:` frontmatter: copy that theme's CSS/JS and use its layout, local or built-in
- [ ] Check the theme exists before wiping the output dir
- [ ] Stop dropping copy errors in `copyDir` and `copyTheme` (use `os.CopyFS`), remove unused `CopyFile`
- [ ] Layouts: page title and description win over site meta, canonical URL is per page

### Markdown
- [ ] Keep raw HTML in `content.md` (`html.WithUnsafe`)
- [ ] Enable GFM: tables, strikethrough, autolinks
- [ ] Auto heading IDs
- [ ] Frontmatter: handle an empty block and a file that opens with a `---` rule

### Library API
- [ ] Let callers silence or redirect `Build` and `Scaffold` output (`io.Writer`)
- [ ] Expose the page list to templates so themes can build a nav
- [ ] `params:` in `ink.yaml` for custom theme data, remove the unused `TemplateData.Theme`
- [ ] Rename before v1: `Page.Content` holds a path, `siteUrl` is the only camelCase key

### CLI and serve
- [ ] Parse flags with the stdlib `flag` package (`--addr=:8080` and unknown flags are ignored today)
- [ ] `serve` prints a broken URL when `--addr` includes a host

### Repo
- [ ] Remove the committed `examples/library/inkssg-library-example` binary and gitignore it
- [ ] Run `go mod tidy` (fsnotify is a direct dependency), check it in CI
- [ ] Tests for frontmatter, `validateOutput`, theme resolution, `injectReloadScript`; run CI with `-race`

## v0.2 — Polish + daily iteration

- [x] `inkssg serve` — local server with file watcher and auto-rebuild
- [x] Copy `static/*` → `public/` (files that need a root URL)
- [ ] Options: `FromConfig`, `WithTheme`, `WithOutputDir`, `WithPagesDir`
- [ ] Hooks: `BeforeBuild`, `AfterBuild`
- [ ] CSS/JS minification
- [ ] Sitemap generation

## v0.3 — Community niceties

- [ ] `inkssg theme eject <name>` — copy built-in theme into local `themes/`
- [ ] `inkssg theme list` — show available themes
- [ ] "Did you mean" suggestions on typos
- [ ] `inkssg validate` — check config without building

## v0.4 — Production niceties

- [ ] `inkssg build --strict` — fail on warnings (broken links, missing alt text)

## v0.5 — Theme sharing

- [ ] `inkssg theme add <git-url>` — clone a theme repo into `themes/`

## Non-goals (forever, unless very strong demand)

- Hot module replacement
- React/JSX / typed templates
- Plugin system
- Theme inheritance
- i18n
- Image optimization pipeline
- Blog post collections / RSS / drafts

If you need those, use Astro, Next, or Hugo.
