[![build status](https://img.shields.io/travis/gcanti/tcomb/master.svg?style=flat-square)](https://travis-ci.org/gcanti/tcomb)
[![dependency status](https://img.shields.io/david/gcanti/tcomb.svg?style=flat-square)](https://david-dm.org/gcanti/tcomb)
![npm downloads](https://img.shields.io/npm/dm/tcomb.svg)

> "Si vis pacem, para bellum" - (Vegetius 5th century)

tcomb is a library for Node.js and the browser which allows you to **check the types** of JavaScript values at runtime with a simple and concise syntax. It's great for **Domain Driven Design** and for adding safety to your internal code.

# TypeScript / Flowtype users

You may want to check out [io-ts](https://github.com/gcanti/io-ts)

# IMPORTANT: Running in production

tcomb is supposed to be used in **development** and is **disabled in production**. 
If you want type checks in production you may use 
- [tcomb-validation](https://github.com/gcanti/tcomb-validation)
- [io-ts](https://github.com/gcanti/io-ts)

# Setup

```sh
npm install tcomb --save
```

**Code example**

A type-checked function:

```js
import t from 'tcomb';

function sum(a, b) {
  t.Number(a);
  t.Number(b);
  return a + b;
}

sum(1, 's'); // throws '[tcomb] Invalid value "s" supplied to Number'

// using babel-plugin-tcomb
function sum(a: number, b: number) {
  return a + b;
}
```

A user defined type:

```js
const Integer = t.refinement(t.Number, (n) => n % 1 === 0, 'Integer');
```

A type-checked class:

```js
const Person = t.struct({
  name: t.String,              // required string
  surname: t.maybe(t.String),  // optional string
  age: t.Integer,                // required integer
  tags: t.list(t.String)       // a list of strings
}, 'Person');

// methods are defined as usual
Person.prototype.getFullName = function () {
  return `${this.name} ${this.surname}`;
};

const person = Person({
  surname: 'Canti'
}); // throws '[tcomb] Invalid value undefined supplied to Person/name: String'
```

Chrome DevTools:

![throws](docs/images/throws.png)

# Docs

- [API](docs/API.md)
- [A little guide to runtime type checking and runtime type introspection](docs/GUIDE.md) (Work in progress)

# Features

**Lightweight**

3KB gzipped, no dependencies.

**Type safety**

All models defined with `tcomb` are type-checked.

**Note**. Instances *are not boxed*, this means that `tcomb` works great with lodash, Ramda, etc. And you can of course use them as props to React components.

**Based on set theory**

- Blog post: [JavaScript, Types and Sets - Part I](https://gcanti.github.io/2014/09/29/javascript-types-and-sets.html)
- Blog post: [JavaScript, Types and Sets - Part II](https://gcanti.github.io/2014/10/07/javascript-types-and-sets-part-II.html)

**Domain Driven Design**

Write complex domain models in a breeze and with a small code footprint. Supported types / combinators:

* user defined types
* structs
* lists
* enums
* refinements
* unions
* intersections
* the option type
* tuples
* dictionaries
* functions
* recursive and mutually recursive types
* interfaces

**Immutability and immutability helpers**

Instances are immutable using `Object.freeze`. This means you can use standard JavaScript objects and arrays. You don't have to change how you normally code. You can update an immutable instance with the provided `update(instance, spec)` function:

```js
const person2 = Person.update(person, {
  name: { $set: 'Guido' }
});
```

where `spec` is an object containing *commands*. The following commands are compatible with the [Facebook Immutability Helpers](http://facebook.github.io/react/docs/update.html):

* `$push`
* `$unshift`
* `$splice`
* `$set`
* `$apply`
* `$merge`

See [Updating immutable instances](docs/API.md#updating-immutable-instances) for details.

**Speed**

`Object.freeze` calls and asserts are executed only in development and stripped out in production (using `process.env.NODE_ENV !== 'production'` tests).

**Runtime type introspection**

All models are inspectable at runtime. You can read and reuse the information stored in your types (in the `meta` static member). See [The meta object](docs/GUIDE.md#the-meta-object) in the docs for details.

Libraries exploiting tcomb's RTI:

- [tcomb-validation](https://github.com/gcanti/tcomb-validation)
- [tcomb-form](https://github.com/gcanti/tcomb-form)
- Blog post: [JSON API Validation In Node.js](https://gcanti.github.io/2014/09/15/json-api-validation-in-node.html)

**Easy JSON serialization / deserialization**

Encodes / decodes your domain models to / from JSON for free.
- Blog post: [JSON Deserialization Into An Object Model](https://gcanti.github.io/2014/09/12/json-deserialization-into-an-object-model.html)

**Debugging with Chrome DevTools**

You can customize the behavior when an assert fails leveraging the power of Chrome DevTools.

```js
// use the default...
t.fail = function fail(message) {
  throw new TypeError('[tcomb] ' + message); // set "Pause on exceptions" on the "Sources" panel for a great DX
};

// .. or define your own behavior
t.fail = function fail(message) {
  console.error(message);
};
```

**Pattern matching**

```js
const result = t.match(1,
  t.String, () => 'a string',
  t.Number, () => 'a number'
);

console.log(result); // => 'a number'
```

**Babel plugin**

Using [babel-plugin-tcomb](https://github.com/gcanti/babel-plugin-tcomb) you can also write (Flow compatible) type annotations:

```js
function sum(a: number, b: number): number {
  return a + b;
}
```

**TypeScript definition file**

[index.d.ts](index.d.ts)

# Contributors

- [Giulio Canti](https://github.com/gcanti) maintainer
- [Becky Conning](https://github.com/beckyconning) `func` combinator ideas and documentation
- [Andrea Lattuada](https://github.com/utaal) `declare` combinator

# How to Build a standalone bundle

```sh
git clone git@github.com:gcanti/tcomb.git
cd tcomb
npm install
npm run dist
```

Will output 2 files:

- `dist/tcomb.js` (development)
- `dist/tcomb.min.js` (production) `Object.freeze` calls and asserts stripped out

# Related libraries

* [tcomb-doc](https://github.com/gcanti/tcomb-doc) Documentation tool for tcomb
* [tcomb-validation](https://github.com/gcanti/tcomb-validation) Validation library based on type combinators
* [tcomb-json-schema](https://github.com/gcanti/tcomb-json-schema) Transforms a JSON Schema to a tcomb type
* [reactuate](https://github.com/reactuate/reactuate) React/Redux stack (not a boilerplate kit)
* [tcomb-react](https://github.com/gcanti/tcomb-react) Alternative syntax for PropTypes
* [mongorito-tcomb](https://github.com/xouabita/mongorito-tcomb) Bring schema validation to Mongorito thanks to tcomb
* [tcomb-form](https://github.com/gcanti/tcomb-form) Forms library for react
* [tcomb-form-types](https://github.com/Industrial/tcomb-form-types) Adds Types/Validations to tcomb-form
* [tcomb-form-native](https://github.com/gcanti/tcomb-form-native) Forms library for react-native
* [tcomb-generate](https://github.com/typeetfunc/tcomb-generate) Generates random data by tcomb type for property-based testing
* [tcomb-additional-types](http://github.com/eserozvataf/tcomb-additional-types) Provides additional predefined types for your existing tcomb setup
* [tcomb-builder](https://github.com/cadre/tcomb-builder) An immutable syntax for defining tcomb types and forms

# Similar projects

* [typed-immutable](https://github.com/Gozala/typed-immutable)
* [immu](https://github.com/scottcorgan/immu)
* [immutable](https://github.com/facebook/immutable-js)
* [mori](https://github.com/swannodette/mori)
* [seamless-immutable](https://github.com/rtfeldman/seamless-immutable)
* [deep-freeze](https://www.npmjs.com/package/deep-freeze)
* [freezer](https://github.com/arqex/freezer)
* [icedam](https://github.com/winkler1/icedam)
* [ObjectModel](https://github.com/sylvainpolletvillard/ObjectModel)
* [rfx](https://github.com/ericelliott/rfx)

# License

The MIT License (MIT)
