# How to contribute to ZITADEL's documentation

Before you start your editing process please take some time to familiarise yourself with the rules for contributing to the documentation of ZITADEL.

## File format

ZITADEL's static site generator uses a Markdown file format. The flavour of the Markdown is specified in the [CAOS Static Site Generator](/docs/cssg.md) document.

## Using headings

Make use of headings to structure and organise the document. The static site generator will make sure that the headings are rendered as a means of navigation for the site's content.

## Headings capitalisation conventions

In order to keep consistency across the documentation we urge you to only write the very first word of a heading with a capital letter. The styling of the static site generator will then guarantee the consistent rendering of first caps, all caps, word caps or whatever style is chosen for the output.

## Using captions

We urge you to always add captions to listings, code snippets, tables, graphics and images, *unless* the image is a logo or icon of any type. Also, inline graphics and images don't require captions due to readability.

Numbering is **required** for captions. See next section.

## Numbering of captions

Numbering has to be continuous using Arabic style numbers. The caption is composed of a prefix and a descriptive text. The prefix is composed of the item type (Listing, Table, Image, Graphic), an incrementing number in Arabic style followed by a period and a description.

If the description consists of **one term** only you must omit any trailing punctuation. If the description represents a complete sentence you must terminate it with a period.

Examples:

* Snippet 3. Example of a protobuf interface deriving ZITADEL's OIDC connector.
* Figure 1. Inside the ESO Atacama Large Millimetre/submillimetre Array.
* Table 21. Representation of all supported quantum-singularity proximation algorithms.
* Image 13. Higgs Boson

## Inline graphics and images

Use sparsely. You should always make sure that inline images do not exceed 48 by 48 pixels (or the equivalent scale for HDPI displays).

## Keyboard shortcuts

When documenting the usage of UI components don't forget to include the keyboard shortcuts available. Use the `kbd:` macro or `<kbd>` HTML tag to style the keys accordingly. To ensure platform-independent documentation always include the modifier-keys for **all** supported platforms, i.e. `[kbd:]Ctrl` or `[kbd:]Opt [kbd:]C`, which renders `<kbd>Ctrl</kbd>` or `<kbd>Opt</kbd> <kbd>C</kbd>`.

## File name conventions

### File name

If and when the file name of a Markdown file consist of more than **one** contiguous word you must use a hyphen `-` to separate the word elements.

Examples:

* content.md
* distribution-guide.md
* programming-examples.md
* explaining-the-code-snippets.md

### File extension

For reasons of consistency the static site generator only accepts Markdown files with the `.md` file extension.

### Internationalisation (I18N)

Writing documentation in country specific languages is highly endorsed. The static site generator makes use of a simple naming convention for Markdown files which are available in different languages:

`{filename}[.{language-id}].md`

Where `filename` represents the name of the file (see section [File name](#file-name)), `language-id` is an *optional* language identifier (see the list of [ISO 639-1 codes](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)), and the mandatory file extension `.md`.

**Note:** If a file name does not specify a language identifier we consider the file name to be implicitly written in English language (defaulting to `.en`). This is also important for the rendering of the language-dependent content navigation as the content navigations always will only include pages in their respective language they represent. This means that Markdown files without language identifier will **only** be rendered in the content navigation of the English documentation variant and will **not be visible at all** in the content navigations of other language variants.
