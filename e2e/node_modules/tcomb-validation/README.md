[![build status](https://img.shields.io/travis/gcanti/tcomb-validation/master.svg?style=flat-square)](https://travis-ci.org/gcanti/tcomb-validation)
[![dependency status](https://img.shields.io/david/gcanti/tcomb-validation.svg?style=flat-square)](https://david-dm.org/gcanti/tcomb-validation)
![npm downloads](https://img.shields.io/npm/dm/tcomb-validation.svg)

A general purpose JavaScript validation library based on type combinators

# Features

- concise yet expressive syntax
- validates native types, refinements, objects, lists and tuples, enums, unions, dicts, intersections
- validates structures with arbitrary level of nesting
- detailed informations on failed validations
- lightweight alternative to JSON Schema
- reuse your domain model written with [tcomb](https://github.com/gcanti/tcomb)

# Documentation

- [Basic usage](#basic-usage)
  - [Primitives](#primitives)
  - [Refinements](#refinements)
  - [Objects](#objects)
  - [Lists and tuples](#lists-and-tuples)
  - [Enums](#enums)
  - [Unions](#unions)
  - [Dicts](#dicts)
  - [Intersections](#intersections)
  - [Nested structures](#nested-structures)
- [Customise error messages](#customise-error-messages)
- [Use cases](#use-cases)
  - [Form validation](#form-validation)
  - [JSON schema](#json-schema)
- [API reference](#api-reference)

# Basic usage

*If you don't know how to define types with tcomb you may want to take a look at its [README](https://github.com/gcanti/tcomb/blob/master/README.md) file.*

The main function is `validate`:

```js
validate(value, type, [options]) -> ValidationResult
```

- `value` the value to validate
- `type` a type defined with the [tcomb](https://github.com/gcanti/tcomb) library
- `options` (optional) is an object with the following keys
  - `path: Array<string | number>` path prefix for validation
  - `context: any` passed to `getValidationErrorMessage` (useful for i18n)
  - `strict: boolean` (default `false`) if `true` no additional properties are allowed while validating structs

returns a `ValidationResult` object containing the result of the validation

**Note**.

- `options` can be an array (as `path` prefix) for backward compatibility (deprecated)

Example

```js
var t = require('tcomb-validation');
var validate = t.validate;

validate(1, t.String).isValid();   // => false
validate('a', t.String).isValid(); // => true
```

You can inspect the result to quickly identify what's wrong:

```js
var result = validate(1, t.String);
result.isValid();             // => false
result.firstError().message;  // => 'Invalid value 1 supplied to String'

// see `result.errors` to inspect all errors
```

## Primitives

```js
// null and undefined
validate('a', t.Nil).isValid();       // => false
validate(null, t.Nil).isValid();      // => true
validate(undefined, t.Nil).isValid(); // => true

// strings
validate(1, t.String).isValid();   // => false
validate('a', t.String).isValid(); // => true

// numbers
validate('a', t.Number).isValid(); // => false
validate(1, t.Number).isValid();   // => true

// booleans
validate(1, t.Boolean).isValid();    // => false
validate(true, t.Boolean).isValid(); // => true

// optional values
validate(null, maybe(t.String)).isValid(); // => true
validate('a', maybe(t.String)).isValid();  // => true
validate(1, maybe(t.String)).isValid();    // => false

// functions
validate(1, t.Function).isValid();              // => false
validate(function () {}, t.Function).isValid(); // => true

// dates
validate(1, t.Date).isValid();           // => false
validate(new Date(), t.Date).isValid();  // => true

// regexps
validate(1, t.RegExp).isValid();    // => false
validate(/^a/, t.RegExp).isValid(); // => true
```

## Refinements

You can express more fine-grained contraints with the `refinement` syntax:

```js
// a predicate is a function with signature: (x) -> boolean
var predicate = function (x) { return x >= 0; };

// a positive number
var Positive = t.refinement(t.Number, predicate);

validate(-1, Positive).isValid(); // => false
validate(1, Positive).isValid();  // => true
```

## Objects

### Structs

```js
// an object with two numerical properties
var Point = t.struct({
  x: t.Number,
  y: t.Number
});

validate(null, Point).isValid();            // => false
validate({x: 0}, Point).isValid();          // => false, y is missing
validate({x: 0, y: 'a'}, Point).isValid();  // => false, y is not a number
validate({x: 0, y: 0}, Point).isValid();    // => true
validate({x: 0, y: 0, z: 0}, Point, { strict: true }).isValid(); // => false, no additional properties are allowed
```

### Interfaces

**Differences from structs**

- also checks prototype keys

```js
var Serializable = t.interface({
  serialize: t.Function
});

validate(new Point(...), Serializable).isValid(); // => false

Point.prototype.serialize = function () { ... }

validate(new Point(...), Serializable).isValid(); // => true
```

## Lists and tuples

**Lists**

```js
// a list of strings
var Words = t.list(t.String);

validate(null, Words).isValid();                  // => false
validate(['hello', 1], Words).isValid();          // => false, [1] is not a string
validate(['hello', 'world'], Words).isValid();    // => true
```

**Tuples**

```js
// a tuple (width x height)
var Size = t.tuple([Positive, Positive]);

validate([1], Size).isValid();      // => false, height missing
validate([1, -1], Size).isValid();  // => false, bad height
validate([1, 2], Size).isValid();   // => true
```

## Enums

```js
var CssTextAlign = t.enums.of('left right center justify');

validate('bottom', CssTextAlign).isValid(); // => false
validate('left', CssTextAlign).isValid();   // => true
```

## Unions

```js
var CssLineHeight = t.union([t.Number, t.String]);

validate(null, CssLineHeight).isValid();    // => false
validate(1.4, CssLineHeight).isValid();     // => true
validate('1.2em', CssLineHeight).isValid(); // => true
```

## Dicts

```js
// a dictionary of numbers
var Country = t.enums.of(['IT', 'US'], 'Country');
var Warranty = t.dict(Country, t.Number, 'Warranty');

validate(null, Warranty).isValid();             // => false
validate({a: 2}, Warranty).isValid();           // => false, ['a'] is not a Country
validate({US: 2, IT: 'a'}, Warranty).isValid(); // => false, ['IT'] is not a number
validate({US: 2, IT: 1}, Warranty).isValid();   // => true
```

## Intersections

```js
var Min = t.refinement(t.String, function (s) { return s.length > 2; }, 'Min');
var Max = t.refinement(t.String, function (s) { return s.length < 5; }, 'Max');
var MinMax = t.intersection([Min, Max], 'MinMax');

MinMax.is('abc'); // => true
MinMax.is('a'); // => false
MinMax.is('abcde'); // => false
```

## Nested structures

You can validate structures with an arbitrary level of nesting:

```js
var Post = t.struct({
  title: t.String,
  content: t.String,
  tags: Words
});

var mypost = {
  title: 'Awesome!',
  content: 'You can validate structures with arbitrary level of nesting',
  tags: ['validation', 1] // <-- ouch!
};

validate(mypost, Post).isValid();             // => false
validate(mypost, Post).firstError().message;  // => 'tags[1] is `1`, should be a `Str`'
```

# Customise error messages

You can customise the validation error message defining a function `getValidationErrorMessage(value, path, context)` on the type constructor:

```js
var ShortString = t.refinement(t.String, function (s) {
  return s.length < 3;
});

ShortString.getValidationErrorMessage = function (value) {
  if (!value) {
    return 'Required';
  }
  if (value.length >= 3) {
    return 'Too long my friend';
  }
};

validate('abc', ShortString).firstError().message; // => 'Too long my friend'
```

## How to keep DRY?

In order to keep the validation logic in one place, one may define a custom combinator:

```js
function mysubtype(type, getValidationErrorMessage, name) {
  var Subtype = t.refinement(type, function (x) {
    return !t.String.is(getValidationErrorMessage(x));
  }, name);
  Subtype.getValidationErrorMessage = getValidationErrorMessage;
  return Subtype;
}

var ShortString = mysubtype(t.String, function (s) {
  if (!s) {
    return 'Required';
  }
  if (s.length >= 3) {
    return 'Too long my friend';
  }
});

```

# Use cases

## Form validation

Let's design the process for a simple sign in form:

```js
var SignInInfo = t.struct({
  username: t.String,
  password: t.String
});

// retrieves values from the UI
var formValues = {
  username: $('#username').val().trim() || null,
  password: $('#password').val().trim() || null
};

// if formValues = {username: null, password: 'password'}
var result = validate(formValues, SignInInfo);
result.isValid();             // => false
result.firstError().message;  // => 'Invalid value null supplied to /username: String'
```

## JSON schema

If you don't want to use a JSON Schema validator or it's not applicable, you can just use this lightweight library in a snap. This is the JSON Schema example of [http://jsonschemalint.com/](http://jsonschemalint.com/)

```json
{
  "type": "object",
  "properties": {
    "foo": {
      "type": "number"
    },
    "bar": {
      "type": "string",
      "enum": [
        "a",
        "b",
        "c"
      ]
    }
  }
}
```

and the equivalent `tcomb-validation` counterpart:

```js
var Schema = t.struct({
  foo: t.Number,
  bar: t.enums.of('a b c')
});
```

let's validate the example JSON:

```js
var json = {
  "foo": "this is a string, not a number",
  "bar": "this is a string that isn't allowed"
};

validate(json, Schema).isValid(); // => false

// the returned errors are:
- Invalid value "this is a string, not a number" supplied to /foo: Number
- Invalid value "this is a string that isn't allowed" supplied to /bar: "a" | "b" | "c"
```

**Note**: A feature missing in standard JSON Schema is the powerful [refinement](#refinements) syntax.

# Api reference

## ValidationResult

`ValidationResult` represents the result of a validation. It containes the following fields:

- `errors`: a list of `ValidationError` if validation fails
- `value`: an instance of `type` if validation succeded

```js
// the definition of `ValidationError`
var ValidationError = t.struct({
  message: t.String,                        // a default message for developers
  actual: t.Any,                            // the actual value being validated
  expected: t.Function,                     // the type expected
  path: list(t.union([t.String, t.Number])) // the path of the value
}, 'ValidationError');

// the definition of `ValidationResult`
var ValidationResult = t.struct({
  errors: list(ValidationError),
  value: t.Any
}, 'ValidationResult');
```

### #isValid()

Returns true if there are no errors.

```js
validate('a', t.String).isValid(); // => true
```

### #firstError()

Returns an object that contains an error message or `null` if validation succeeded.

```js
validate(1, t.String).firstError().message; // => 'value is `1`, should be a `Str`'
```

## validate(value, type, [options]) -> ValidationResult

- `value` the value to validate
- `type` a type defined with the tcomb library
- `options` (optional) is an object with the following keys
  - `path: Array<string | number>` path prefix for validation
  - `context: any` passed to `getValidationErrorMessage` (useful for i18n)
  - `strict: boolean` (default `false`) if `true` no additional properties are allowed while validating structs

# Tests

Run `npm test`

# License

The MIT License (MIT)
