---
title: Doc structure
description: This document explains how syntax is rendered.
---

### General Structure

To deploy a successful build, you need to have a `docs` folder on your root directory.

This folder contains subfolders that are mapped to routes of your doc page later,
a `index.svelte` file for the homepage, and
a `static` folder for all assets and metadata.

The `index.svelte` as well as the static folder with a `manifest.json` is mandatory, otherwise the build will fail.
Other files referenced by the homepage or within markdown has to be included to the static folder.

```bash
├ docs
│ ├ get_started
│ │ ├ 00-indroduction.en.md
│ │ ├ 00-indroduction.de.md
│ │ ├ 00-indroduction.it.md
│ │ ├ 01-get-started.en.md
│ │ ├ 01-get-started.de.md
│ │ ├ 01-get-started.it.md
│ │ ├ 02-concluding.en.md
│ │ ├ 02-concluding.de.md
│ │ └ 02-concluding.it.md
│ ├ api
│ ├ ├ 00-indroduction.en.md
│ │ ├ 01-get-started.en.md
│ │ └ 02-concluding.en.md
  ├ static
  │ ├ manifest.json
  │ ├ favicon.ico
  │ ├ android-chrome-192x192.png
  └ index.svelte
```

According to the configuration above, two routes with the names `getting_started` and `api` are generated.

A doc page consists of one ore more `markdown` files, fetched in sorting order. 
So prefixing with numbers might be a good choice for organizing your structure.

Anchors are automatically generated using headings in the documentation and by default they are latinised.

### Markdown Files

Markdown files contain a meta data section including the `title`, used as section header later and other possibly helpful information about the author of the page

---

Take a look at this example:

```markdown
<!-- start of file -->
---
title: Introduction
---
<!-- markdown here -->
```

Markdown files are compiled to html and styled with a custom theme scheme (currently a work in progress).

Currently the `marked.js` renderer is taking care of this markdown files, treating some syntax as custom and some as default.
Take a look at the sections below to get an understanding of whats going on.

#### links

Links to external targets (starting with http) are opened in a new tab.

If you want to reference to an other article of your doc page, you can use the `title` concatenated to your `route`.

---
This link takes you to [General Structure](get_started#General_Structure).

```md
[General Structure](get_started#General_Structure)
```

#### headings

Headings of level three and four are slugged, so headings like
```md
### header
#### header-2 
```

will appear in the navigation and can be referenced by links.

> The meta-data title and headings of your markdown files are slugged, so that a navigation can be build.
> Make sure all titles and headings are distinct, otherwise your build will end up failing!

#### code

Code is highlighted accoding to the specified language

```md
    ```yaml
    config: main.yml
    ```
```

#### code switcher

Code is highlighted accoding to the specified language

```html
    <a href="link"></a>
```

```js
    console.log('');
```

```yaml
    config: main.yml
```

#### hr

hr can be used to split content, but it will work for code blocks on the right only.
```md
---

Left block (can by anything)

Right block (only code)
```

will be rendered to

---

Left block can be <strong>anything</strong>.

```js
    console.log('This can be any code block');
```
This can be useful to explain a code snipped.
#### other

##### blockquote
if you have to highlight something noteworthy.

```md
> This is a blockquote and it contains important infos
```

> This is a blockquote and it contains important infos

##### html

You can also use raw HTML in your Markdown, and it'll mostly work pretty well.

```html
<div style="border: 1px solid white; padding: 1rem; font-style: italic;">hello</div>
```

<div style="border: 1px solid white; padding: 1rem; font-style: italic; margin-bottom: 20px;">hello</div>

##### other markdown syntax

These are other blocks which are highlighted but have standard conventions.

* list
* listitem
* paragraph
* table
* tablerow
* tablecell

### Translating the API docs

Default language is set in manifest.json. Make sure an attribute `lang` is set.

Anchors are automatically generated using headings in the documentation and by default they are latinised to make sure the URL is always conforming to RFC3986.

If we need to translate the API documentation to a language using unicode chars, we can setup this app to export the correct anchors by setting up `SLUG_PRESERVE_UNICODE` to `true` in `config.js`.