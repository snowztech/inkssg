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

Detects `pages/`, themes, and assets automatically.

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

Put images, fonts, and other static files in `assets/`. They copy to `public/assets/` as-is.

```
my-site/
└── assets/
    └── img/
        └── logo.png
```

Reference them with absolute paths in your content:

```markdown
![logo](/assets/img/logo.png)
```

## Static files

Anything in `static/` is copied to the root of `public/` as-is. Use it for files that need a root URL: `robots.txt`, a `CNAME`, a standalone HTML page, an image you want at `/photo.jpg`.

```
my-site/
└── static/
    ├── robots.txt       → /robots.txt
    └── demo.html        → /demo.html
```

If a static file has the same name as a built page, the page wins.

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
- `examples/multi-page/` — multiple pages, shared theme