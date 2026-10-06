# inkssg

Static site generator for simple sites. No config needed, no 20 dependencies.

```
go get github.com/snowztech/inkssg
```

## Quick start

Create a site and build it with 3 lines of Go:

```go
package main

import inkssg "github.com/snowztech/inkssg"

func main() {
    inkssg.Build()
}
```

Run it:

```
go run main.go
```

That's it. inkssg detects your site structure, builds pages to `public/`, and outputs a summary.

## Library usage

### Build from current directory

```go
inkssg.Build()
```

Detects `pages/`, themes, `assets/` and `static/` automatically.

### Build from a specific directory

```go
inkssg.Build("./my-site")
```

Same behavior, explicit path.

### Error handling

```go
if err := inkssg.Build(); err != nil {
    log.Fatal(err)
}
```

## CLI usage

Install:

```
go install github.com/snowztech/inkssg/cmd/inkssg@latest
```

Commands:

```
inkssg build      build current directory
inkssg new <path> scaffold a new site
inkssg version   show version
```

## Site structure

A minimal inkssg site:

```
my-site/
├── pages/
│   └── index/
│       └── content.html
└── assets/
    └── image.png
```

Or with markdown:

```
my-site/
├── pages/
│   └── index/
│       └── content.md
└── assets/
    └── image.png
```

No `ink.yaml` needed for simple sites.

Everything inkssg looks for, all of it optional except `pages/`:

```
my-site/
├── ink.yaml       site-wide data
├── pages/         one folder per page          → /<name>.html
├── themes/        your own themes
├── assets/        images, fonts, icons         → /assets/...
└── static/        files that need a root URL   → /...
```

## Pages

Each page lives in `pages/<name>/`. The folder name becomes the output filename.

```
pages/
├── index/content.html    → index.html
├── about/content.md      → about.html
└── contact/content.md  → contact.html
```

Page content goes in `content.md` or `content.html`. Both work.

## Frontmatter

Add metadata at the top of any content file:

```markdown
---
title: About
description: A short blurb
theme: minimal
---
```

Supported fields:

- `title` — page title
- `description` — meta description
- `theme` — theme name (defaults to `minimal`)

## Themes

Themes live in `themes/<name>/`. Only `layout.html` is required.

```
themes/
└── minimal/
    ├── layout.html   (required)
    ├── styles.css  (optional)
    └── script.js  (optional)
```

### Built-in themes

`minimal` is built in and loaded automatically. No setup needed.

### Custom themes

Create your own by adding a `themes/` folder:

```
my-site/
├── pages/
│   └── index/
│       └── content.html
└── themes/
    └── landing/
        └── layout.html
```

Use `theme: landing` in frontmatter to apply it.

### Layout template

Your `layout.html` receives these data objects:

```html
<!DOCTYPE html>
<html>
<head>
  <title>{{.Page.Title}}</title>
  <meta name="description" content="{{.Page.Description}}">
</head>
<body>
  {{.Content}}
</body>
</html>
```

| Object | What it contains |
|--------|-----------------|
| `.Content` | rendered page content |
| `.Page` | title, description from frontmatter |
| `.Site` | site data from `ink.yaml` (only if it exists) |
| `.Theme` | theme data |

## Config (optional)

Add `ink.yaml` when you need site-wide data:

```yaml
name: My Site
default_theme: landing
output_dir: dist
```

Supported fields:

- `name` — site name (available as `.Site.Name` in templates)
- `default_theme` — theme to use if not set in frontmatter
- `output_dir` — output directory (defaults to `public/`). Must be inside the project root.
- `pages_dir` — where pages live (defaults to `pages/`)

## ink.yaml vs frontmatter

Use frontmatter for page-specific data (title, description, theme).

Use `ink.yaml` for site-wide data you want available on every page (name, links, bio, meta).

Example frontmatter:

```markdown
---
title: About
description: A short blurb
---
```

Example `ink.yaml`:

```yaml
name: Vikusha
meta:
  lang: en
  description: Go framework for AI assistants
links:
  - name: GitHub
    url: https://github.com/snowztech/vikusha
```

`.Site` is only populated when `ink.yaml` exists. If there's no `ink.yaml`, templates work fine without it.

## Assets

Put the files your pages use in `assets/`: images, fonts, icons. They copy to `public/assets/` as-is.

```
my-site/
└── assets/
    └── img/
        └── logo.png     → /assets/img/logo.png
```

Reference them with absolute paths in your content:

```markdown
![logo](/assets/img/logo.png)
```

## Static files

Some files need an exact URL at the root of the site. Put those in `static/`. Everything in it copies to the root of `public/` as-is, keeping its folders.

```
my-site/
└── static/
    ├── robots.txt                  → /robots.txt
    ├── CNAME                       → /CNAME
    ├── demo.html                   → /demo.html
    └── .well-known/
        └── security.txt            → /.well-known/security.txt
```

The folder is optional and needs no config. If it isn't there, nothing happens.

Typical uses:

- `robots.txt`, `ads.txt`, `CNAME`, `.nojekyll`
- site verification files from search consoles
- a standalone HTML page you wrote by hand and don't want wrapped in a theme
- a file you want to link to directly, like `/cv.pdf`

### assets/ or static/?

| | `assets/` | `static/` |
|---|---|---|
| Served at | `/assets/<path>` | `/<path>` |
| Use for | files your pages and themes reference | files that need an exact root URL |
| Example | `assets/img/logo.png` → `/assets/img/logo.png` | `static/robots.txt` → `/robots.txt` |

If you're not sure, use `assets/`. Reach for `static/` only when the URL matters.

### What inkssg does not do

Static files are copied, not built. No markdown rendering, no layout, no frontmatter. An HTML file in `static/` is published exactly as you wrote it. If you want a page with your theme around it, make it a page under `pages/`.

Dotfiles are copied too, so `.well-known/` and `.nojekyll` work. So does anything your OS drops in there, like `.DS_Store`.

### Name clashes

The build copies in this order, and later steps overwrite earlier ones:

1. theme files → `public/themes/<name>/`
2. `assets/` → `public/assets/`
3. `static/` → `public/`
4. pages → `public/<name>.html`

So a built page always wins: `static/index.html` does not replace the page from `pages/index/`. A file under `static/assets/` or `static/themes/` does replace the one from `assets/` or the theme. Avoid those two folder names inside `static/` unless that is what you want.

### With `inkssg serve`

`static/` is watched like the rest of the project. Add or change a file and the site rebuilds.

## Output

Build output:

```
✓ index → public/index.html
✓ about → public/about.html
built 2 pages in 89ms
```

On error:

```
✗ contact: missing content.md or content.html
1 page failed (of 3) in 12ms
```

The `public/` folder is ready to deploy.

## Examples

See `examples/` in the repo for complete sites:

- `examples/library/` — minimal site using the library API
- `examples/minimal/` — single page with default theme
- `examples/multi-page/` — multiple pages, shared theme, a `robots.txt` in `static/`