# Cypress terminal report [![Build Status](https://circleci.com/gh/archfz/cypress-terminal-report/tree/master.svg?style=svg)](https://app.circleci.com/pipelines/github/archfz/cypress-terminal-report) [![Downloads](https://badgen.net/npm/dw/cypress-terminal-report)](https://www.npmjs.com/package/cypress-terminal-report) [![Version](https://badgen.net/npm/v/cypress-terminal-report)](https://www.npmjs.com/package/cypress-terminal-report)

> This documentation is for cypress >= 10.0.0. For older versions see [3.x.x branch](https://github.com/archfz/cypress-terminal-report/tree/3.x.x).

<div align="center">

[Limitations](#limitations-and-notes)
• [Install](#install)
• [Options](#options)
• [Logging after/before all](#logging-after-all-and-before-all-hooks)
• [Logging to files](#logging-to-files)
• [Development](#development)
• [Release Notes](#release-notes)

</div>

Plugin for cypress that adds better terminal output for easier debugging. 
Prints cy commands, browser console logs, cy.request and cy.intercept data. Great for your pipelines.

* looks pretty in console
* logs all commands, requests and browser console logs
* supports logging to files
* option between logging only on failure (default) or always
* options for trimming and compacting logs
* support for multiple and nested mocha contexts
* log commands from before all and after all hooks ([with a catch*](#logging-after-all-and-before-all-hooks))

Try it out by cloning [cypress-terminal-report-demo](https://github.com/archfz/cypress-terminal-report-demo).


![demo](https://raw.githubusercontent.com/archfz/cypress-terminal-report/master/demo.png?sanitize=true#gh-light-mode-only)
![demo](https://raw.githubusercontent.com/archfz/cypress-terminal-report/master/demo_dark.png?sanitize=true#gh-dark-mode-only)


## Limitations and notes

- By default logs are not printed for successful tests. Please see [option](#optionsprintlogstoconsole) to change this.
- `console.log` usage was never meant to be used in the cypress test code. Using it will
  not log anything with this plugin. Using it also goes against the queue nature of 
  cypress. Use `cy.log` instead. [See here for more details](https://github.com/archfz/cypress-terminal-report/issues/67).
- Using `cypress-fail-fast` and logging to files does not work out of the box. To enable support use 
the [`logToFilesOnAfterRun`](#optionslogtofilesonafterrun) option.

## Install

### Requirements

- `>=4.0.0` requires cypress `>=10.0.0`
- `>=3.0.0` requires cypress `>=4.10.0`
- `<3.0.0` requires cypress `>=3.8.0`

1. Install npm package.
    ```bash
    npm i --save-dev cypress-terminal-report
    ```
2. If using typescript and es6 imports ensure `esModuleInterop` is enabled.
3. Register the output plugin in `cypress.config.{js|ts}`
    ```js
    module.exports = defineConfig({
      e2e: {
        setupNodeEvents(on, config) {
          require('cypress-terminal-report/src/installLogsPrinter')(on);
        }
      }
    });
    ```
4. Register the log collector support in `cypress/support/e2e.{js|ts}`
    ```js
    require('cypress-terminal-report/src/installLogsCollector')();
    ```

## Options

<br/>

### _Options for the plugin install_

> require('cypress-terminal-report/src/installLogsPrinter')(on, options)

#### `options.defaultTrimLength`
integer; default: 800; Max length of cy.log and console.warn/console.error.

#### `options.commandTrimLength` 
integer; default: 800; Max length of cy commands.

#### `options.routeTrimLength`
integer; default: 5000; Max length of cy.route, cy.request or XHR data.

#### `options.compactLogs` 
integer?; default: null; If it is set to a number greater or equal to 0, this amount of logs 
will be printed only around failing commands. Use this to have shorter output especially 
for when there are a lot of commands in tests. When used with `options.printLogsToConsole=always` 
for tests that don't have any `severity=error` logs nothing will be printed.

#### `options.outputCompactLogs` 
integer? | false; default: null; Overrides `options.compactLogs` for the file log output specifically, 
when `options.outputTarget` is specified. Allows compacting of the terminal and the file output logs to different levels.  
If `options.outputCompactLogs` is unspecified, file output will use `options.compactLogs`.
If set to `false`, output file logs will not compact even if `options.compactLogs` is set.

#### `options.outputRoot` 
string; default: null; Required if `options.outputTarget` provided. [More details](#logging-to-files).

#### `options.specRoot` 
string; default: null; Cypress specs root relative to package json. [More details](#log-specs-in-separate-files ).

#### `options.outputTarget`
object; default: null; Output logs to files. [More details](#logging-to-files).

#### `options.outputVerbose`
boolean; default: true; Toggles verbose output.

#### `options.printLogsToConsole`
string; Default: 'onFail'. When to print logs to console, possible values: 'always', 'onFail', 'never' - When set to always
logs will be printed to console for successful tests as well as failing ones.

#### `options.printLogsToFile`
string; Default: 'onFail'. When to print logs to file(s), possible values: 'always', 'onFail', 'never' - When set to always
logs will be printed to file(s) for successful tests as well as failing ones.

#### `options.includeSuccessfulHookLogs`
boolean; Default: false. Commands from before all and after all hooks by default get logged only if
a command from them failed. This default is in accordance with the defaults on `options.printLogsTo*` to 
avoid printing too many, possibly irrelevant, information. However you can set this to `true` if you 
need more extensive logging, but be aware that commands will be logged to terminal from these hooks 
regardless whether there were failing tests in the suite. This is because we can't know for sure in 
advanced if a test fails or not.

#### `options.collectTestLogs` *1
([spec, test, state], [type, message, severity][]) => void; default: undefined;
Callback to collect each test case's logs after its run.
The first argument contains information about the test: the `spec` (test file), `test` (test title) and `state` (test state) fields.
The second argument contains the test logs. 'type' is from the same list as for the `collectTypes` support install option (see below). Severity can be of ['', 'error', 'warning'].

#### `options.logToFilesOnAfterRun`
boolean; default: false;
When set to true it enables additional log write pass to files using the cypress [`after:run`](https://docs.cypress.io/api/plugins/after-run-api) plugin
hook. This option can only be used with cypress 6.2.0 onwards, and with the additional 
`experimentalRunEvents` configuration on versions smaller than 6.7.0.


<br/>

### _Options for the support install_

> require('cypress-terminal-report/src/installLogsCollector')(options);

#### `options.collectTypes` 
array; default: ['cons:log','cons:info', 'cons:warn', 'cons:error', 'cy:log', 'cy:xhr', 'cy:request', 'cy:route', 'cy:intercept', 'cy:command']
What types of logs to collect and print. By default all types are enabled. The 'cy:command' is the general type that
contain all types of commands that are not specially treated.

#### `options.filterLog` 
null | ([type, message, severity]) => boolean; default: undefined; 
Callback to filter logs manually.
The type is from the same list as for the `collectTypes` option. Severity can be of ['', 'error', 'warning'].

#### `options.processLog`
null | ([type, message, severity]) => string; default: undefined;
Callback to process logs manually.
The type is from the same list as for the `collectTypes` option. Severity can be of ['', 'error', 'warning'].

#### `options.collectTestLogs` *2
(mochaRunnable, [type, message, severity][]) => void; default: undefined;
Callback to collect each test case's logs after its run.
The `mochaRunnable` is of type `Test | Hook` from the mocha library.
The type is from the same list as for the `collectTypes` option. Severity can be of ['', 'error', 'warning'].

#### `options.xhr.printHeaderData` 
boolean; default false; Whether to print header data for XHR requests.

#### `options.xhr.printRequestData`
boolean; default false; Whether to print request data for XHR requests besides response data.

#### `options.enableExtendedCollector`
boolean; default false; Enables an extended collector which will also collect command logs from
before all and after all hooks.

#### `options.enableContinuousLogging`
boolean; default false; Enables logging logs to terminal continuously / immediately as they are registered.
This feature is unstable and has an impact on pipeline performance. This option has no effect for extended 
collector, only works for the simple collector. Use only for debugging purposes in case the pipelines / 
tests are timing out. 

> NOTE: In case of this option enabled, logs will come before the actual title of the test. Also the 
> `printLogsToConsole` option will be ignored. Logging to files might also get impacted.

#### _Example for options for the support install_

```js
  // ...
  // Options for log collector
  const options = {
    // Log console output only
    collectTypes: ["cons:log", "cons:info", "cons:warn", "cons:error"],
  };

  // Register the log collector
  require("cypress-terminal-report/src/installLogsCollector")(options);
  // ...
```

## Logging after all and before all hooks

Commands from before all and after all hooks are not logged by default. A new experimental feature introduces
support for logging commands from these hooks: [`enableExtendedCollector`](#optionsenableextendedcollector).
This feature is by default disabled as it relies much more heavily on internals of cypress and 
mocha, thus **there is a higher chance of something breaking, especially with cypress upgrades**.

Once the feature enabled, logs from these hooks will only appear in console if:
- a command from the hook fails
- hook passes and [`printLogsToConsole`](#optionsprintlogstoconsole) == `always`
- hook passes and [`printLogsToConsole`](#optionsprintlogstoconsole) ==`onFail` 
  and [`includeSuccessfulHookLogs`](#optionsprintlogstoconsole) == `true`


## Logging to files

To enable logging to file you must add the following configuration options to the
plugin install.

```js
setupNodeEvents(on, config) {
  // ...
  const options = {
    outputRoot: config.projectRoot + '/logs/',
    outputTarget: {
      'out.txt': 'txt',
      'out.json': 'json',
    }
  };

  require('cypress-terminal-report/src/installLogsPrinter')(on, options);
  // ...
}
```

The `outputTarget` needs to be an object where the key is the relative path of the
file from `outputRoot` and the value is the __type__ of format to output. 

Supported types: `txt`, `json`.

### Log specs in separate files

To create log output files per spec file instead of one single file change the
key in the `outputTarget` to the format `{directory}|{extension}`, where
`{directory}` the root directory where to generate the files and `{extension}`
is the file extension for the log files. The generated output will have the
same structure as in the cypress specs root directory.

```js
setupNodeEvents(on, config) {
  const options = {
    outputRoot: config.projectRoot + '/logs/',
    // Used to trim the base path of specs and reduce nesting in the generated output directory.
    specRoot: 'cypress/e2e',
    outputTarget: {
      'cypress-logs|json': 'json',
    }
  };
}
```

### Custom output log processor

If you need to output in a custom format you can pass a function instead of a string
to the `outputTarget` value. This function will be called with the list of messages
per spec per test. It is called right after one spec finishes, which means on each
iteration it will receive for one spec the messages. See for example below.

> NOTE: The chunks have to be written in a way that after every write the file is 
in a valid format. This has to be like this since we cannot detect when cypress
runs the last test. This way we also make the process faster because otherwise the
more tests would execute the more RAM and processor time it would take to rewrite 
all the logs to the file.

Inside the function you will have access to the following API:

- `this.size` - Current char size of the output file.
- `this.atChunk` - The count of the chunk to be written.
- `this.initialContent` - The initial content of the file. Defaults to ''. Set this
before the first chunk write in order for it to work.
- `this.chunkSeparator` - Chunk separator string. Defaults to ''. This string will 
be written between each chunk. If you need a special separator between chunks use this 
as it is internally handled to properly write and replace the chunks.
- `this.writeSpecChunk(specPath, dataString, positionInFile?)` - Writes a chunk of 
data in the output file.

```js
  // ...
  const options = {
    outputTarget: {
      'custom.output': function (allMessages) {
        // allMessages= {[specPath: string]: {[testTitle: string]: [type: string, message: string, severity: string][]}}

        Object.entries(allMessages).forEach(([spec, tests]) => {
          let text = `${spec}:\n`
          Object.entries(tests).forEach(([test, messages]) => {
            text += `    ${test}\n`
            messages.forEach(([type, message, severity]) => {
              text += `        ${type} (${severity}): ${message}\n`
            })
          })
          
          // .. Process the tests object into desired format ..
          // Insert chunk into file, by default at the end.
          this.writeSpecChunk(spec, text);
          // Or before the last two characters.
          this.writeSpecChunk(spec, text, -2);
        });
      }
    }
  };
  // ...
```

See [JsonOutputProcessor](./src/outputProcessor/JsonOutputProcessor.js) implementation as a
good example demonstrating both conversion of data into string and chunk write position
alternation.

## Development

### Testing

Tests can be found under `/test`. The primary expectations are run with mocha and these tests in fact
start cypress run instances and assert on their output. So that means there is a cypress suite that
is used to emulate the usage of the plugin, and a mocha suite to assert on those emulations.

To add tests you need to first add a case to existing cypress spec or create a new one and then
add the case as well in the `/test/test.js`. To run the tests you can use `npm test` in the test \
directory. You should add `it.only` to the test case you are working on to speed up development.

## Release Notes

#### 4.1.1

- Fix issue with `http` module causing `vite` issues. [issue](https://github.com/archfz/cypress-terminal-report/issues/154) by [wopian](https://github.com/wopian) 
- Dependency updates. by [wopian](https://github.com/wopian)

#### 4.1.0

- Add experimental [`enableContinuousLogging`](#optionsenablecontinuouslogging) option for timeout debugging purposes. [issue](https://github.com/archfz/cypress-terminal-report/issues/150)

#### 4.0.3

- Fix issue with errors throw outside of tests overlapping real error. [issue](https://github.com/archfz/cypress-terminal-report/issues/152)
- Add additional potential source for spec file path determination.
- Update cypress in tests to 10.3.0 to confirm support.

#### 4.0.2

- Typescript typing fix to support both esm and commonjs require. [issue](https://github.com/archfz/cypress-terminal-report/issues/151)

#### 4.0.1

- Typescript typing fix and readme update. [issue](https://github.com/archfz/cypress-terminal-report/issues/148)

#### 4.0.0

- Add support for cypress ^10. [Follow cypress upgrade](https://deploy-preview-4186--cypress-docs.netlify.app/guides/references/migration-guide#Migrating-to-Cypress-version-10-0).
- ! Breaking change: `specRoot` option cannot be calculated anymore using config, as 
  `integrationFolder` option was removed in cypress. [This now has to be set manually](#log-specs-in-separate-files). 

#### 3.5.2

- Fix issue where top-level `.spec` files that call test functions in other files results in multiple output files being created.  by [bvandercar-vt](https://github.com/bvandercar-vt)
- Security dependency updates.

#### 3.5.1

- Fix custom output processor example in README. by [bvandercar-vt](https://github.com/bvandercar-vt)
- Add more exported types to facilitate creating custom output processors in TypeScript. by [bvandercar-vt](https://github.com/bvandercar-vt)
- Security dependency updates.

#### 3.5.0

- Add new feature [`outputCompactLogs`](#optionsoutputcompactlogs) to allow for optionally overriding compactLogs for the output file specifically. by [bvandercar-vt](https://github.com/bvandercar-vt) [issue](https://github.com/archfz/cypress-terminal-report/issues/138)
- Fix typing for `processLog`. [issue](https://github.com/archfz/cypress-terminal-report/issues/132)
- Security dependency updates.
- Update cypress to 9.5.x in tests to confirm support.

#### 3.4.2

- Fix incorrectly typed message type arguments. [issue](https://github.com/archfz/cypress-terminal-report/issues/132)
- Security updates.
- Update cypress to 9.4.1 in tests to confirm support.

#### 3.4.1

- Add severity typescript types. [merge-request](https://github.com/archfz/cypress-terminal-report/pull/128)  by [tebeco](https://github.com/tebeco)
- Fix nested file logging by fallback for determining spec path. [issue](https://github.com/archfz/cypress-terminal-report/issues/124)
- Update cypress to 9.1.0 in tests to confirm support.

#### 3.4.0

- Implement `cy:fetch` logging. [merge-request](https://github.com/archfz/cypress-terminal-report/pull/126) by [Erik-Outreach](https://github.com/Erik-Outreach)
- Update cypress to 8.7.0 in tests to confirm support.

#### 3.3.4

- Fix issue in `extedend control` where skip tests would double consume logs, and cause domain exception. [issue](https://github.com/archfz/cypress-terminal-report/issues/122)
- Update cypress to 8.6.0 in tests to confirm support.

#### 3.3.3

- Fix issue with `cy.intercept` overrides not working. [issue](https://github.com/archfz/cypress-terminal-report/issues/120)
- Update cypress to 8.5.0 in tests to confirm support.

#### 3.3.2

- Fix issue with no response on XHR breaking tests. [issue](https://github.com/archfz/cypress-terminal-report/issues/119)

#### 3.3.1

- Fix issue `cy:intercept` not between the allowed configuration options. [issue](https://github.com/archfz/cypress-terminal-report/issues/113)
- Fix issue with plugin breaking cypress with skipped tests. [issue1](https://github.com/archfz/cypress-terminal-report/issues/116) [issue2](https://github.com/archfz/cypress-terminal-report/issues/114)
- Update cypress to 8.3.0 in tests to confirm support.

#### 3.3.0

- Added support for logging command from skipped tests. [issue](https://github.com/archfz/cypress-terminal-report/issues/111)
- Update cypress to 8.1.0 in tests to confirm support.

#### 3.2.2

- Added protection against incorrect tabbing level determined from test parents breaking logging to terminal. [issue](https://github.com/archfz/cypress-terminal-report/issues/107) 
- Remove peer dependency mocha.
- Update cypress to 7.4.0 in tests to confirm support.

#### 3.2.1

- Additional fix over extended control with nested mocha contexts and after each hooks failing. [issue](https://github.com/archfz/cypress-terminal-report/issues/98)

#### 3.2.0

- Introduce support for [cypress-fail-fast](https://github.com/javierbrea/cypress-fail-fast) with the [logToFilesOnAfterRun](#optionslogtofilesonafterrun) option. [issue](https://github.com/archfz/cypress-terminal-report/issues/91)

#### 3.1.2

- Fix issue with duplicated log send on extended control when parent suite has after each but current suite doesn't. [issue](https://github.com/archfz/cypress-terminal-report/issues/98)
- Fix issue with empty array of logs causing extra unnecessary new lines in console output. [issue](https://github.com/archfz/cypress-terminal-report/issues/99)
- Confirm support with cypress 7.3.0.

#### 3.1.1

- Fix compatibility with cypress < 6.0.0. [issue](https://github.com/archfz/cypress-terminal-report/issues/95)
- Confirm support with cypress 7.2.0.

#### 3.1.0

- Add support for `cy.intercept()` capturing and logging. [issue](https://github.com/archfz/cypress-terminal-report/issues/87)
- Improve `cy.request()` timeout and xhr abort error logging. [issue](https://github.com/archfz/cypress-terminal-report/issues/89) 
- Update cypress in tests to 7.1.0 to confirm support.

#### 3.0.4

- Fix skipped root suites breaking extended collector. [merge-request](https://github.com/archfz/cypress-terminal-report/pull/90) by [jmoses](https://github.com/jmoses)

#### 3.0.3

- Fix incomplete previous release.

#### 3.0.2

- Fixed issue with errors not displaying correctly for commands outside of tests. [issue](https://github.com/archfz/cypress-terminal-report/issues/85)
- Update readme with additional notes on limitations.

#### 3.0.1

- Fix issue with cucumber tests not logging properly to nested files. [issue](https://github.com/archfz/cypress-terminal-report/issues/82)
- Fix issue with `filterLog` and `processLog` options running too soon on non-final log list. [issue](https://github.com/archfz/cypress-terminal-report/issues/84)

#### 3.0.0

- ! Breaking change in [`options.collectTestLogs`](#optionscollecttestlogs-2). First parameter (previously called context) changed.
- ! Possibly breaking change: Test names in output files now contain mocha contexts.
- ! Minimum required cypress version increased to 4.10.0.
- ! Removed deprecated `printLogs` option on support install.
- Added support for [logging commands from before all and after all hooks](#logging-after-all-and-before-all-hooks). [issue](https://github.com/archfz/cypress-terminal-report/issues/55)
- Added support for multiple and nested mocha context. In console logs are tabbed according to nesting
level of context and in output files context titles are always added. [issue](https://github.com/archfz/cypress-terminal-report/issues/70)
- Added support for [option to process logs](#optionsprocesslog). [issue](https://github.com/archfz/cypress-terminal-report/issues/74) [merge-request](https://github.com/archfz/cypress-terminal-report/pull/81) by [FLevent29](https://github.com/FLevent29)
- Added support to [disable verbose logging](#optionsoutputverbose). [issue](https://github.com/archfz/cypress-terminal-report/issues/71) [merge-request](https://github.com/archfz/cypress-terminal-report/pull/80) by [FLevent29](https://github.com/FLevent29)
- Updated unicode icons for log types on linux.
- Updated cypress to 6.6.x in tests to confirm support.

#### 2.4.0

- Improved logging of xhr with status code, response data in case of failure and duration. [issue](https://github.com/archfz/cypress-terminal-report/issues/61) [merge-request](https://github.com/archfz/cypress-terminal-report/pull/66) by [peruukki](https://github.com/peruukki)   
- Added `console.debug` logging support. [merge-request](https://github.com/archfz/cypress-terminal-report/pull/48) by [reynoldsdj](https://github.com/reynoldsdj)
- Improve config schema for trimming options to allow `number`.

#### 2.3.1

- Fixed issue in nested file output for spec file names containing multiple dots. [merge-request](https://github.com/rinero-bigpanda) by [rinero-bigpanda](https://github.com/rinero-bigpanda)
- Fixed issue in nested file output not cleaning up existing files. [issue](https://github.com/archfz/cypress-terminal-report/issues/57) 
- Update to cypress 6.0.0 in tests and fix expectations.

#### 2.3.0

- Added support for custom log collector function on both [nodejs](#optionscollecttestlogs-1) and [browser](#optionscollecttestlogs-2) sides. [issue](https://github.com/archfz/cypress-terminal-report/issues/56) [merge-request](https://github.com/archfz/cypress-terminal-report/pull/59) by [peruukki](https://github.com/peruukki)
- Improved cy.request logging when log set to false and when there are network errors. [merge-request](https://github.com/archfz/cypress-terminal-report/pull/60) by [bjowes](https://github.com/bjowes) 

#### 2.2.0

- Added support for [logging each spec to its own file](#log-specs-in-separate-files). [issue](https://github.com/archfz/cypress-terminal-report/issues/49)

#### 2.1.0

- ! Separated option `printLogs` to [`printLogsToConsole`](#optionsprintlogstoconsole) and [`printLogsToFile`](#optionsprintlogstofile).
`printLogs` won't work anymore and will print a warning in the cypress logs. Read documentation on how to
upgrade. [issue](https://github.com/archfz/cypress-terminal-report/issues/47) [merge-request](https://github.com/archfz/cypress-terminal-report/pull/51) by [FLevent29](https://github.com/FLevent29)
- Typescript typings support added. [issue](https://github.com/archfz/cypress-terminal-report/issues/37) [merge-request](https://github.com/archfz/cypress-terminal-report/pull/52) by [bengry](https://github.com/bengry)

#### 2.0.0

- Removed deprecated exports from index.js. If you were still using require from index.js please
see [installation](#install) for updating.
- Added JSON schema validation for options to prevent invalid options and assumptions. [issue](https://github.com/archfz/cypress-terminal-report/issues/45)
- Fixed issue where output to file would insert at incorrect position for JSON when ran from GUI.
- Reworked the file output processing code and thus the API changed as well. Custom output processors
will have to be updated to current API when upgrading to this version. Check [readme section](#custom-output-log-processor).
- Added printing to terminal of time spent in milliseconds for output files to be written.
- Improved Error instanceof checking for console log arguments printing. [issue](https://github.com/archfz/cypress-terminal-report/issues/41)
- Update cypress to 5.0.0 in tests to confirm compatibility.

#### 1.4.2

- Fixed issue with compact logs breaking after each hook with message `undefined is not iterable`. [issue](https://github.com/archfz/cypress-terminal-report/issues/39)
- Update cypress to 4.11.0 in tests to confirm compatibility.

#### 1.4.1

- Added log to terminal of the list of generated output files. [merge-request](https://github.com/archfz/cypress-terminal-report/issues/29) by [@bjowes](https://github.com/bjowes)
- Fixed line endings for output on windows. [merge-request](https://github.com/archfz/cypress-terminal-report/issues/29) by [@bjowes](https://github.com/bjowes)
- Fixed incorrect readme with inverse installation requires for plugin and support. [merge-request](https://github.com/archfz/cypress-terminal-report/pull/36) by [@zentby](https://github.com/zentby) 

#### 1.4.0

- Added new feature to compact output of logs, [see here](#optionscompactlogs). [issue](https://github.com/archfz/cypress-terminal-report/issues/27)
- Fixed incorrect severity on cons:error and cons:warn.
- Fixed compatibility with cypress 4.8. [issue-1](https://github.com/archfz/cypress-terminal-report/issues/35) 
[issue-2](https://github.com/archfz/cypress-terminal-report/issues/34) 
[issue-3](https://github.com/archfz/cypress-terminal-report/issues/33)
- Fixed issue with webpack compatibility caused by native includes getting in compilation files. For this please revise
the installation documentation and change the requires for the install of this plugin. Deprecated require by index. [issue](https://github.com/archfz/cypress-terminal-report/issues/32)
- Fixed issue with logsChainId not being reset and causing test failures and potentially live failures with error
'Cannot set property '2' of undefined'. [issue](https://github.com/archfz/cypress-terminal-report/issues/31)

#### 1.3.1

- Added creation of output path for outputTarget files if directory does not exist.
- Fixed issue with ctrAfter task not being registered when outputTarget not configured. [issue](https://github.com/archfz/cypress-terminal-report/issues/26)

#### 1.3.0

- Added support for [logging to file](#writing-logs-to-files), with builtin support for json and text and possible custom processor. [issue](https://github.com/archfz/cypress-terminal-report/issues/17) 
- Added support for logging XHR request body and also headers for requests and responses. [issue](https://github.com/archfz/cypress-terminal-report/issues/25)
- Reformatted the log message for route and request commands.
- Replace all tab characters with spaces on console log.

#### 1.2.1

- Fixed issue with incorrect command being marked as failing when there are additional logs after the actual failing one. [issue](https://github.com/archfz/cypress-terminal-report/issues/24)
- Fixed issue where console would receive undefined and the plugin would break: split on undefined. [issue](https://github.com/archfz/cypress-terminal-report/issues/23)
- Bumping default trim lengths for cy:command and cons:* log types.
- Improvements on logging objects from console. 
- Fix incorrect severity calculation for cy:route.
- Fix yellow output on powershell.
- Fix windows console icons with different set. [issue](https://github.com/archfz/cypress-terminal-report/issues/7)
- Update cypress version for testing to 4.3.0.
- Set peer version of cypress to >=3.8.1. [issue](https://github.com/archfz/cypress-terminal-report/issues/22)

#### 1.2.0

- Fixed issue with cy.request accepting parameters in multiple formats and the plugin not recognizing this. [merge-request](https://github.com/archfz/cypress-terminal-report/pull/19) by [@andrew-blomquist-6](https://github.com/andrew-blomquist-6)
- Improved browser console logs for Error and other objects. [issue-1](https://github.com/archfz/cypress-terminal-report/issues/18) [issue-2](https://github.com/archfz/cypress-terminal-report/issues/16)
- Added support for filtering logs. See `collectTypes` and `filterLog` options for the support install. from [issue](https://github.com/archfz/cypress-terminal-report/issues/15).
- Removed option `printConsoleInfo` in favor of above. Also now the console.log and info are by
 default enabled.

#### 1.1.0

- Added notice for logs not appearing in dashboard. from [issue](https://github.com/archfz/cypress-terminal-report/issues/8)
- Added [support for logging](#options) console.info and console.log. in [issue](https://github.com/archfz/cypress-terminal-report/issues/12) 
- Added better logging of cy.requests. in [issue](https://github.com/archfz/cypress-terminal-report/issues/12) by [@zhex900](https://github.com/zhex900)

#### 1.0.0

- Added tests and CI to repository.
- Added support for showing logs even for successful tests. in [issue](https://github.com/archfz/cypress-terminal-report/issues/3) by [@zhex900](https://github.com/zhex900)
- Fixed issue with incorrectly labeled failed commands. in [issue](https://github.com/archfz/cypress-terminal-report/issues/3) by [@zhex900](https://github.com/zhex900)
- Fixed issue with logs from cy.route breaking tests on XHR API type of requests. [merge-request](https://github.com/archfz/cypress-terminal-report/pull/1) by [@zhex900](https://github.com/zhex900)
