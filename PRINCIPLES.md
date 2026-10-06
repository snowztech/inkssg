# Principles

Small sites should not need 20 dependencies and a webpack config.

## Why inkssg exists

Sometimes you just need a one-page site. You shouldn't need a framework to build it.

inkssg is for when you want to edit one file, run one command, and deploy. No config needed for simple sites.

## Principles

**1. Library-first, CLI second**

The primary way to use inkssg is as a Go library. `go get github.com/snowztech/inkssg`, write 3 lines, build. CLI is a wrapper around the library.

**2. Simpler over more powerful**

Every feature added is a burden for every user. inkssg does the common case well, not every case. If you need more, use Hugo.

**3. Explicit over magical**

No hidden rewrites, no implicit conventions beyond the essentials. Paths are literal: `assets/<path>` is served at `/assets/<path>`, `static/<path>` at `/<path>`. No surprise path rewriting or auto-redirects.

**4. Small footprint, not comprehensive**

No hot reload, no plugins, no i18n, no blog features. A small binary you can read in an afternoon. Less code means fewer bugs and easier contributions.

**5. Convention over configuration**

Structure is predictable. `pages/`, `ink.yaml`, `themes/`. If you know one inkssg site, you know them all.

## Target audience

- Landing pages
- Bio or link pages
- Small documentation
- Personal sites

Not for: blogs, large documentation sites, complex web apps.

## When to use

| Use this | Not this |
|---------|----------|
| inkssg: simple site, few pages, Go user | Hugo: feature-rich, needs templating knowledge |
| Astro: React-heavy, complex interactivity | Jekyll: Ruby, plugin-heavy |

## Non-goals

inkssg will never have:

- Hot module replacement
- Plugin system
- Theme inheritance
- i18n
- Blog features (posts, drafts, RSS)
- Image optimization pipeline

If you need those, Hugo or Astro are better tools.

## Decisions

Technical decisions are in `ignored/DECISIONS.md` (internal).

## Docs

Full user documentation is in [docs/inkssg.md](docs/inkssg.md).