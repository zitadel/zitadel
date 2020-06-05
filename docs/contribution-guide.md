![ZITADEL](/docs/img/zitadel-logo-oneline-lightdesign@2x.png "Zitadel Logo")

# How to contribute to Zitadel's documentation
Before you start your editing process please take some time to familiarise yourself with the rules for contributing to the documentation of Zitadel.

## Using captions
We urge you to always add captions to listings, code snippets, tables, graphics and images, *unless* the image is a logo or icon of any type. Also, inline graphics and images don't require captions due to readability.

Numbering is **required** for captions. See next section.

## Numbering of captions
Numbering has to be continuous using arabic style numbers. The caption is composed of a prefix and a descriptive text. The prefix is composed of the item type (Listing, Table, Image, Graphic), and incrementing number in arabic style followed by a period and a description.

If the description consists of **one term** only you must omit any trailing punctuation. If the description represents a complete sentense you must terminate it with a period.

Examples:

* Snippet 3. Example of a protobuf interface deriving ZITADEL's OIDC connector.
* Figure 1. Inside the ESO Atacama Large Millimeter/submillimeter Array.
* Table 21. Representation of all supported quantum-singularity proximation algorithms.
* Image 13. Higgs Boson

## Inline graphics and images
Use sparsely. You should always make sure that inline images do not exceed 48 by 48 pixels (or the equivalent scale for HDPI displays).

## Keyboard shortcuts
When documenting the usage of UI components don't forget to include the keyboard shortcuts available. Use the `kbd:` macro to style the keys accordingly. To ensure platform-independent documentation always include the modifier-keys for **all** supported platforms, i.e. [kbd:]Ctrl or [kbd:]Opt [kbd:]C.
