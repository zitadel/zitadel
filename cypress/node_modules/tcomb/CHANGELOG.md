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

# v3.2.29

- **Bug Fix**
  - [typescript] fix `interface`'s `extend` signature, #329 (@apepper)

# v3.2.28

- **Bug Fix**
  - Enums.is() with an array value should always be false, #327 (@phorsuedzie)

# v3.2.27

- **Bug Fix**
  - struct typescript extend signature, #317 (@lramel)
  - defaultProps missing from StructOptions, #317 (@lramel)

# v3.2.25

- **Bug Fix**
  - fromJSON makes use of defaultProps while deserializing a struct, fix #312 (@rkmax)

# v3.2.24

- **Bug Fix**
  - Struct extension of a refinement of a struct now use the correct displayName, fix #297 (@gcanti)

# v3.2.23

- **Bug Fix**
  - declare: remove unnecessary limitation, fix #291 (@gcanti)

# v3.2.22

- **Polish**
  - Update TypeScript definitions to allow module augmentation (@RedRoserade)

# v3.2.21

- **Bug Fix**
  - TypeScript definition file: `Nil` should be `void | null` (@francescogior)

# v3.2.20

- **Polish**
  - add `options` (struct, interface) to typescript definition (@gcanti)

# v3.2.19

- **Polish**
  - add `strict` (struct, interface) to typescript definition (@gcanti)

# v3.2.18

- **Bug Fix**
  - fix `define` in typescript definition (@gcanti)

# v3.2.17

- **Bug Fix**
  - add missing `t.Integer` to typescript definition (@gcanti)

# v3.2.16

- **Bug Fix**
  - strict structs with additional methods should not throw on updating, fix #267 (@gcanti)

# v3.2.15

- **New Feature**
  - Added support for overwriting `defaultProps` in `t.struct.extend`, fix #257 (@tehnomaag)

# v3.2.14

- **Bug Fix**
  - replace `instanceof Array` with `Array.isArray`, fix #255 (@ewnd9)

# v3.2.13

- **Bug Fix**
  - fromJSON: typecasting of values inside `t.intersection`, fix #250 (@gcanti)

# v3.2.12

- **Bug Fix**
  - now `interface` doesn't filter additional props when props contain a struct, fix #245 (@gcanti)

# v3.2.11

- **Bug Fix**
  - allow declare'd unions with custom dispatch, fix #242 (@gcanti)

# v3.2.10

- **Bug Fix**
  - handle nully values in interface `is` function (@gcanti)

# v3.2.9

- **New Feature**
  - fromJSON: track error path, fix #235 (@gcanti)
- **Internal**
  - change shallow copy in order to improve perfs (@gcanti)

# v3.2.8

- **Bug Fix**
  - mixing types and classes in a union throws, fix #232 (@gcanti)

# v3.2.7

- **Bug Fix**
  - add support for class constructors, `fromJSON` module (@gcanti)
  - type-check the value returned by a custom reviver, `fromJSON` module (@gcanti)

# v3.2.6

- **Bug Fix**
  - null Maybes should stringify to null, fix #227 (@gcanti)

# v3.2.5

- **Polish**
  - prevent bugs when enums are defined through `t.declare` (@gcanti)

# v3.2.4

- **Polish**
  - decouple usage of new operator in create() function, fix #223 (@gcanti)

# v3.2.3

- **Polish**
  - add `isNil` check in interface constructor
- **Experimental**
  - add support for [babel-plugin-tcomb](https://github.com/gcanti/babel-plugin-tcomb), fix #218 (@gcanti)

# v3.2.2

- **Bug Fix**
  - relax `isObject` contraint in interface combinator, fix #214

# v3.2.1

- **Bug Fix**
  - fix missing path argument in FuncType
- **Polish**
  - better stringify serialization for functions

# v3.2.0

- **New Feature**
  - `isSubsetOf` module, function for determining whether one type is compatible with another type (@R3D4C73D)
  - default props for structs (thanks @timoxley)
- **Documentation**
  - global strict settings are deprecated (see https://github.com/gcanti/tcomb/issues/168#issuecomment-222422999)

# v3.1.0

- **New Feature**
  - add `t.Integer` to standard types
  - add `t.Type` to standard types
  - `interface` combinator, fix #195, [docs](https://github.com/gcanti/tcomb/blob/master/docs/API.md#the-interface-combinator) (thanks @ctrlplusb)
    - add interface support to fromJSON (@minedeljkovic)
  - add support for extending refinements, fix #179, [docs](https://github.com/gcanti/tcomb/blob/master/docs/API.md#extending-structs)
  - local and global `strict` option for structs and interfaces, fix #203, [docs](https://github.com/gcanti/tcomb/blob/master/docs/API.md#strictness)
  - Chrome Dev Tools custom formatter for tcomb types [docs](https://github.com/gcanti/tcomb/blob/master/docs/API.md#the-libinstalltypeformatter-module)
- **Bug Fix**
  - More intelligent immutability update handling, fix #199 (thanks @ctrlplusb)
  - func combinator: support optional arguments, fix #198 (thanks @ivan-kleshnin)
- **Internal**
  - add "Struct" prefix to structs default name
  - `mixin()` now allows identical references for overlapping properties

# v3.0.0

**Warning**. If you don't rely in your codebase on the property `maybe(MyType)(undefined) === null` this **is not a breaking change** for you.

- **Breaking Change**
  - prevent `Maybe` constructor from altering the value when `Nil`, fix #183 (thanks @gabro)

# v2.7.0

- **New Feature**
  - `lib/fromJSON` module: generic deserialize, fix #169
  - `lib/fromJSON` TypeScript definition file
- **Bug Fix**
  - t.update module: $apply doesn't play well with dates and regexps, fix #172
  - t.update: cannot $merge and $remove at once, fix #170 (thanks @grahamlyus)
  - TypeScript: fix Exported external package typings file '...' is not a module
  - misleading error message in `Struct.extend` functions, fix #177 (thanks @Firfi)

# v2.6.0

- **New Feature**
  - `declare` API: recursive and mutually recursive types (thanks @utaal)
  - typescript definition file, fix #160 (thanks @DanielRosenwasser)
  - `t.struct.extend`, fix #164 (thanks @dzdrazil)
- **Internal**
  - split main file to separate modules, fix #158
  - add "typings" field to package.json (TypeScript)
  - add `predicate` field to irreducibles meta objects
- **Documentation**
  - revamp [API.md](https://github.com/gcanti/tcomb/blob/master/docs/API.md)
  - add ["A little guide to runtime type checking and runtime type introspection"](https://github.com/gcanti/tcomb/blob/master/docs/GUIDE.md) (WIP)

## v2.5.2

- **Bug Fix**
  - remove the assert checking if the type returned by a union dispatch function is correct (was causing issues with unions of unions or unions of intersections)

## v2.5.1

- **Internal**
  - `t.update` should not change the reference when no changes occur, fix #153

# v2.5.0

- **New Feature**
  - check if the type returned by a union dispatch function is correct, fix #136 (thanks @fcracker79)
  - added `refinement` alias to `subtype` (which is deprecated), fix #140
- **Internal**
  - optimisations: for identity types return early in production, fix #135 (thanks @fcracker79)
  - exposed `getDefaultName` on combinator constructors

## v2.4.1

- **New Feature**
  - added struct multiple inheritance, fix #143

# v2.4.0

- **New Feature**
  - unions
    - added `update` function, #127
    - the default `dispatch` implementation now handles unions of unions, #126
    - show the offended union type in error messages

# v2.3.0

- **New Feature**
  - Add support for lazy messages in asserts, fix #124
  - Better error messages for assert failures, fix #120

  The messages now have the following general form:

  ```
  Invalid value <value> supplied to <context>
  ```

  where context is a slash-separated string with the following properties:

  - the first element is the name of the "root"
  - the following elements have the form: `<field name>: <field type>`

  Note: for more readable messages remember to give types a name

  Example:

  ```js
  var Person = t.struct({
    name: t.String
  }, 'Person'); // <- remember to give types a name

  var User = t.struct({
    email: t.String,
    profile: Person
  }, 'User');

  var mynumber = t.Number('a');
  // => Invalid value "a" supplied to Number

  var myuser = User({ email: 1 });
  // => Invalid value 1 supplied to User/email: String

  myuser = User({ email: 'email', profile: { name: 2 } });
  // => Invalid value 2 supplied to User/profile: Person/name: String
  ```


## v2.2.1

- **Experimental**
  - pattern matching #121

# v2.2.0

- **New Feature**
  - added `intersection` combinator fix #111

    **Example**

    ```js
    const Min = t.subtype(t.String, function (s) { return s.length > 2; }, 'Min');
    const Max = t.subtype(t.String, function (s) { return s.length < 5; }, 'Max');
    const MinMax = t.intersection([Min, Max], 'MinMax');

    MinMax.is('abc'); // => true
    MinMax.is('a'); // => false
    MinMax.is('abcde'); // => false
    ```

- **Internal**
  - optimised the generation of default names for types

# v2.1.0

- **New Feature**
  - added aliases for pre-defined irreducible types fix #112
  - added overridable `stringify` function to handle error messages and improve performances in development (replaces the experimental `options.verbose`)

## v2.0.1

- **Experimental**
  - added `options.verbose` (default `true`) to handle messages (set `options.verbose = false` to improve performances in development)

# v2.0.0

- **New Feature**
  - add support to types defined as ES6 classes #99
  - optimized for production code: asserts and freeze only in development mode
  - add `is(x, type)` function
  - add `isType(x)` function
  - add `stringify(x)` function
- **Breaking change**
  - numeric types on enums #93  (thanks @m0x72)
  - remove asserts when process.env.NODE_ENV === 'production' #100
  - do not freeze if process.env.NODE_ENV === 'production' #103
  - func without currying #96 (thanks @tmcw)
  - remove useless exports #104
  - drop bower support #101
  - remove useless exports
    * Type
    * slice
    * shallowCopy
    * getFunctionName
