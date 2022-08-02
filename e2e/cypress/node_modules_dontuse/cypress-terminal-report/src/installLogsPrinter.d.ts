/// <reference types="cypress" />

declare function installLogsPrinter(on: Cypress.PluginEvents, options?: installLogsPrinter.PluginOptions): void;
declare namespace installLogsPrinter {

  type Severity = '' | 'error' | 'warning';

  type LogType = 'cons:log' |
    'cons:info' |
    'cons:warn' |
    'cons:error' |
    'cons:debug' |
    'cy:log' |
    'cy:xhr' |
    'cy:fetch' |
    'cy:request' |
    'cy:route' |
    'cy:intercept' |
    'cy:command' |
    'ctr:info';

  type AllMessages = {
    [specPath: string]: {
      [testTitle: string]: [type: LogType, message: string, severity: Severity][]
    }
  };

  interface PluginOptions {
    /**
     * Max length of cy.log and console.warn/console.error.
     * @default 800
     */
    defaultTrimLength?: number;

    /**
     * Max length of cy commands.
     * @default 800
     */
    commandTrimLength?: number;

    /**
     * Max length of cy.route request data.
     * @default 5000
     */
    routeTrimLength?: number;

    /**
     * If it is set to a number greater or equal to 0, this amount of logs will be printed only around failing commands. Use this to have shorter output especially for when there are a lot of commands in tests. When used with options.printLogs=always for tests that don't have any severity=error logs nothing will be printed.
     * @default null
     */
    compactLogs?: number | null;

    /**
     * If it is set to a number greater or equal to 0, will override compactLogs for the file log output specifically. Use this for compacting of the terminal and the file output logs to different levels.
     * @default null
     */
    outputCompactLogs?: false | number | null;

    /**
     * Required if outputTarget provided. [More details](https://github.com/archfz/cypress-terminal-report#logging-to-files).
     * @default null
     */
    outputRoot?: string | null;

    /**
     * Output logs to files. [More details](https://github.com/archfz/cypress-terminal-report#logging-to-files).
     * @default null
     */
    outputTarget?: Record<
      string,
      | 'json'
      | 'txt'
      | ((allMessages: AllMessages) => string)
      >;

    /**
     * Toggles verbose output.
     * @default true
     */
    outputVerbose?: boolean | true;

    /**
     * Cypress specs root relative to package json. [More details](https://github.com/archfz/cypress-terminal-report#logging-to-files).
     * @default null
     */
    specRoot?: string | null;

    /**
     * When set to always logs will be printed for console for successful test as well as failing ones.
     * @default 'onFail'
     */
    printLogsToConsole?: 'onFail' | 'always' | 'never';

    /**
     * When set to always logs will be printed to file for successful test as well as failing ones.
     * @default 'onFail'
     */
    printLogsToFile?: 'onFail' | 'always' | 'never';

    /**
     * Whether to log commands from hooks that passed.
     * If enabled even when all tests pass in a spec the commands will always
     * be printed from before and after hooks.
     * @default false
     */
    includeSuccessfulHookLogs: boolean

    /**
     * Callback to collect each test case's logs after its run.
     * @default undefined
     */
    collectTestLogs?: (context: {spec: string, test: string, state: string}, messages: [/* type: */ LogType, /* message: */ string, /* severity: */ Severity][]) => void;
  }
}
export = installLogsPrinter;
