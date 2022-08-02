# mochawesome

[![npm](https://img.shields.io/npm/v/mochawesome.svg?style=flat-square)](http://www.npmjs.com/package/mochawesome) ![Node.js CI](https://github.com/adamgruber/mochawesome/workflows/Node.js%20CI/badge.svg) [![Gitter](https://img.shields.io/gitter/room/nwjs/nw.js.svg?style=flat-square)](https://gitter.im/mochawesome/general)

Mochawesome is a custom reporter for use with the Javascript testing framework, [mocha][mocha]. It runs on Node.js (>=10) and works in conjunction with [mochawesome-report-generator][marge] to generate a standalone HTML/CSS report to help visualize your test runs.

## Features

<img align="right" src="./docs/marge-report-1.0.1.png" alt="Mochawesome Report" width="55%" />

- Simple, clean, and modern design
- Beautiful charts (via ChartJS)
- Support for test and suite nesting
- Displays before and after hooks
- Review test code inline
- Stack trace for failed tests
- Support for adding context information to tests
- Filters to display only the tests you want
- Responsive and mobile-friendly
- Offline viewing
- Supports `parallel` mode

## Usage

1. Add Mochawesome to your project:

`npm install --save-dev mochawesome`

2. Tell mocha to use the Mochawesome reporter:

`mocha testfile.js --reporter mochawesome`

3. If using mocha programatically:

```js
var mocha = new Mocha({
  reporter: 'mochawesome',
});
```

### Parallel Mode

Since `mocha@8` test files can be run in parallel using the `--parallel` flag. In order for mochawesome to work properly it needs to be registered as a hook.

`mocha tests --reporter mochawesome --require mochawesome/register`

### Output

Mochawesome generates the following inside your project directory:

```
mochawesome-report/
├── assets
│   ├── app.css
│   ├── app.js
│   ├── MaterialIcons-Regular.woff
│   ├── MaterialIcons-Regular.woff2
│   ├── roboto-light-webfont.woff
│   ├── roboto-light-webfont.woff2
│   ├── roboto-medium-webfont.woff
│   ├── roboto-medium-webfont.woff2
│   ├── roboto-regular-webfont.woff
│   └── roboto-regular-webfont.woff2
├── mochawesome.html
└── mochawesome.json
```

The two main files to be aware of are:

**mochawesome.html** - The rendered report file

**mochawesome.json** - The raw json output used to render the report

### Options

Options can be passed to the reporter in two ways.

#### Environment variables

The reporter will try to read environment variables that begin with `MOCHAWESOME_`.

```bash
$ export MOCHAWESOME_REPORTFILENAME=customReportFilename
```

_Note that environment variables must be in uppercase._

#### Mocha reporter-options

You can pass comma-separated options to the reporter via mocha's `--reporter-options` flag. Options passed this way will take precedence over environment variables.

```bash
$ mocha test.js --reporter mochawesome --reporter-options reportDir=customReportDir,reportFilename=customReportFilename
```

Alternately, `reporter-options` can be passed in programatically:

```js
var mocha = new Mocha({
  reporter: 'mochawesome',
  reporterOptions: {
    reportFilename: 'customReportFilename',
    quiet: true,
  },
});
```

#### Available Options

The options below are specific to the reporter. For a list of all available options see [mochawesome-report-generator options][marge-options].

| Option Name       | Type    | Default     | Description                                                                                                                       |
| :---------------- | :------ | :---------- | :-------------------------------------------------------------------------------------------------------------------------------- |
| `quiet`           | boolean | false       | Silence console messages                                                                                                          |
| `reportFilename`  | string  | mochawesome | Filename of saved report (html and json) <br> _See [notes](#reportfilename-replacement-tokens) for available token replacements._ |
| `html`            | boolean | true        | Save the HTML output for the test run                                                                                             |
| `json`            | boolean | true        | Save the JSON output for the test run                                                                                             |
| `consoleReporter` | string  | spec        | Name of mocha reporter to use for console output, or `none` to disable console report output entirely                             |

#### reportFilename replacement tokens

Using the following tokens it is possible to dynamically alter the filename of the generated report.

- **[name]** will be replaced with the spec filename when possible.
- **[status]** will be replaced with the status (pass/fail) of the test run.
- **[datetime]** will be replaced with a timestamp. The format can be - specified using the `timestamp` option.

For example, given the spec `cypress/integration/sample.spec.js` and the following config:

```
{
  reporter: "mochawesome",
  reporterOptions: {
    reportFilename: "[status]_[datetime]-[name]-report",
    timestamp: "longDate"
  }
}
```

The resulting report file will be named `pass_February_23_2022-sample-report.html`

**Note:** The `[name]` replacement only occurs when mocha is running one spec file per process and outputting a separate report for each spec. The most common use-case is with Cypress.

### Adding Test Context

Mochawesome ships with an `addContext` helper method that can be used to associate additional information with a test. This information will then be displayed inside the report.

**Please note: arrow functions will not work with `addContext`.** See the [example](#example).

### `addContext(testObj, context)`

| param   | type           | description                         |
| :------ | :------------- | :---------------------------------- |
| testObj | object         | The test object                     |
| context | string\|object | The context to be added to the test |

**Context as a string**

Simple strings will be displayed as is. If you pass a URL, the reporter will attempt to turn it into a link. If the URL links to an image or video, it will be shown inline.

**Context as an object**

Context passed as an object must adhere to the following shape:

```js
{
  title: 'some title'; // must be a string
  value: {
  } // can be anything
}
```

#### Example

Be sure to use ES5 functions and not ES6 arrow functions when using `addContext` to ensure `this` references the test object.

```js
const addContext = require('mochawesome/addContext');

describe('test suite', function () {
  it('should add context', function () {
    // context can be a simple string
    addContext(this, 'simple string');

    // context can be a url and the report will create a link
    addContext(this, 'http://www.url.com/pathname');

    // context can be an image url and the report will show it inline
    addContext(this, 'http://www.url.com/screenshot-maybe.jpg');

    // context can be an object with title and value properties
    addContext(this, {
      title: 'expected output',
      value: {
        a: 1,
        b: '2',
        c: 'd',
      },
    });
  });
});
```

It is also possible to use `addContext` from within a `beforeEach` or `afterEach` test hook.

```js
describe('test suite', () => {
  beforeEach(function () {
    addContext(this, 'some context');
  });

  afterEach(function () {
    addContext(this, {
      title: 'afterEach context',
      value: { a: 1 },
    });
  });

  it('should display with beforeEach and afterEach context', () => {
    // assert something
  });
});
```

## Typescript

This project does not maintain its own type definitions, however they are available on npm from [DefinitelyTyped](https://github.com/DefinitelyTyped/DefinitelyTyped/tree/master/types/mochawesome).

```
$ npm install --save-dev @types/mochawesome
```

## Related

[mochawesome-report-generator][marge]

## License

mochawesome is [MIT licensed][license].

[mocha]: https://mochajs.org/
[marge]: https://github.com/adamgruber/mochawesome-report-generator
[marge-options]: https://github.com/adamgruber/mochawesome-report-generator#options
[changelog]: CHANGELOG.md
[license]: LICENSE.md
