# Changelog

> **Tags:**
> - [New Feature]
> - [Bug Fix]
> - [Breaking Change]
> - [Documentation]
> - [Internal]
> - [Polish]
> - [Experimental]

**Note**: Gaps between patch versions are faulty/broken releases.
**Note**: A feature tagged as Experimental is in a high state of flux, you're at risk of it changing without notice.

## v3.4.1

- **Bug Fix**
  - TypeScript: fix recursive type definition error (@RedRoserade)

## v3.4.0

- **New Feature**
  - Add TypeScript definitions (@RedRoserade)

## v3.3.0

- **New Feature**
  - ValidationError: relax message type (from `String` to `Any`), fix #54 (@pigoz)

## v3.2.2

- **Bug Fix**
  - maybe value is not converted from `undefined` to `null` after validation, fix #44 (@gcanti)

## v3.2.1

- **Bug Fix**
  - No error returned by intersection of structs, fix #42 (thanks @chrisui)

## v3.2.0

- **New Feature**
  - add support for default props (tcomb v3.2)

## v3.1.0

- **New Feature**
  - add support for `interface` combinator (tcomb v3.1)

## v3.0.0

**Warning**. If you don't rely in your codebase on the property `maybe(MyType)(undefined) === null` this **is not a breaking change** for you.

- **Breaking Change**
  - upgrade to `tcomb` v3.0.0

## v2.3.0

- **New Feature**
  - add `strict` option: no additional properties are allowed while validating structs, fix #12

## v2.2.0

- **New Feature**
  - replaced `path` argument with `options`, fix #27 (thanks @th0r)

## v2.1.1

- **Experimental**
  - Add support for custom error messages #23

## v2.1.0

- **New Feature**
  - upgrade to tcomb v2.2
    - validation of `intersection` types
    - error messages are formatted accordingly

## v2.0.1

- **Internal**
  - upgrade to tcomb v2.1
  - updated tests accordingly

## v2.0.0

- **Breaking change**
  - upgrade to tcomb v2.0 #17
  - drop bower support #18

## v1.0.4

- **Bug Fix**
    + Path in dict errors includes all keys upto the error #15

## v1.0.3

- **New Feature**
    + Add an optional `path` param to `validate` in order to provide a prefix to error paths #13

## v1.0.2

- **Internal**
    + react-native compatibility #11

## v1.0.1

- **Internal**
    + Move tcomb to peerDependencies #10

## v1.0

Initial release

