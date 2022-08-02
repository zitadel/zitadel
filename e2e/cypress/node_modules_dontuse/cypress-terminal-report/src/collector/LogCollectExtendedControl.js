const CONSTANTS = require('../constants');
const LogCollectBaseControl = require('./LogCollectBaseControl');

/**
 * Collects and dispatches all logs from all tests and hooks.
 */
module.exports = class LogCollectExtendedControl extends LogCollectBaseControl {

  constructor(collectorState, config) {
    super();
    this.config = config;
    this.collectorState = collectorState;

    this.registerCypressBeforeMochaHooksSealEvent();
  }

  register() {
    this.collectorState.setStrict(true);

    this.registerState();
    this.registerBeforeAllHooks();
    this.registerAfterAllHooks();
    this.registerTests();
    this.registerLogToFiles();
  }

  sendLogsToPrinter(logStackIndex, mochaRunnable, options = {}) {
    if (!mochaRunnable.parent.invocationDetails && !mochaRunnable.invocationDetails) {
      return;
    }

    let testState = options.state || mochaRunnable.state;
    let testTitle = options.title || mochaRunnable.title;
    let testLevel = 0;

    let spec = this.getSpecFilePath(mochaRunnable);
    let wait = typeof options.wait === 'number' ? options.wait : 6;

    {
      let parent = mochaRunnable.parent;
      while (parent && parent.title) {
        testTitle = `${parent.title} -> ${testTitle}`
        parent = parent.parent;
        ++testLevel;
      }
    }

    const prepareLogs = () => {
      return this.prepareLogs(logStackIndex, {mochaRunnable, testState, testTitle, testLevel});
    };

    if (options.noQueue) {
      this.debounceNextMochaSuite(Promise.resolve()
        // Need to wait for command log update debounce.
        .then(() => new Promise(resolve => setTimeout(resolve, wait)))
        .then(() => {
          Cypress.backend('task', {
            task: CONSTANTS.TASK_NAME,
            arg: {
              spec: spec,
              test: testTitle,
              messages: prepareLogs(),
              state: testState,
              level: testLevel,
              consoleTitle: options.consoleTitle,
              isHook: options.isHook,
              continuous: false,
            }
          })
            // For some reason cypress throws empty error although the task indeed works.
            .catch((error) => {/* noop */})
        }).catch(console.error)
      );
    } else {
      // Need to wait for command log update debounce.
      cy.wait(wait, {log: false})
        .then(() => {
          cy.task(
            CONSTANTS.TASK_NAME,
            {
              spec: spec,
              test: testTitle,
              messages: prepareLogs(),
              state: testState,
              level: testLevel,
              consoleTitle: options.consoleTitle,
              isHook: options.isHook,
              continuous: false,
            },
            {log: false}
          );
        });
    }
  }

  registerState() {
    const self = this;

    Cypress.on('log:changed', (options) => {
      if (options.state === 'failed') {
        this.collectorState.updateLogStatusForChainId(options.id);
      }
    });
    Cypress.mocha.getRunner().on('test', (test) => {
      this.collectorState.startTest(test);
    });
    Cypress.mocha.getRunner().on('suite', () => {
      this.collectorState.startSuite();
    });
    Cypress.mocha.getRunner().on('suite end', () => {
      this.collectorState.endSuite();
    });

    // Keeps track of before and after all hook indexes.
    Cypress.mocha.getRunner().on('hook', function (hook) {
      if (!hook._ctr_hook && !hook.fn._ctr_hook) {
        // After each hooks get merged with the test.
        if (hook.hookName !== "after each") {
          self.collectorState.addNewLogStack();
        }

        // Before each hooks also get merged with the test.
        if (hook.hookName === "before each") {
          self.collectorState.markCurrentStackFromBeforeEach();
        }

        if (hook.hookName === "before all") {
          self.collectorState.incrementBeforeHookIndex();
        } else if (hook.hookName === "after all") {
          self.collectorState.incrementAfterHookIndex();
        }
      }
    });
  }

  registerBeforeAllHooks() {
    const self = this;

    // Logs commands from before all hook if the hook passed.
    Cypress.mocha.getRunner().on('hook end', function (hook) {
      if (hook.hookName === "before all" && self.collectorState.hasLogsCurrentStack() && !hook._ctr_hook) {
        self.debugLog('extended: sending logs of passed after all hook');
        self.sendLogsToPrinter(
          self.collectorState.getCurrentLogStackIndex(),
          this.currentRunnable,
          {
            state: 'passed',
            isHook: true,
            title: self.collectorState.getBeforeHookTestTile(),
            consoleTitle: self.collectorState.getBeforeHookTestTile(),
          }
        );
      }
    });

    // Logs commands from before all hooks that failed.
    Cypress.on('before:mocha:hooks:seal', function () {
      self.prependBeforeAllHookInAllSuites(this.mocha.getRootSuite().suites, function ctrAfterAllPerSuite() {
        if (
          this.test.parent === this.currentTest.parent // Since we have after all in each suite we need this for nested suites case.
          && this.currentTest.failedFromHookId // This is how we know a hook failed the suite.
          && self.collectorState.hasLogsCurrentStack()
        ) {
          self.debugLog('extended: sending logs of failed before all hook');
          self.sendLogsToPrinter(
            self.collectorState.getCurrentLogStackIndex(),
            this.currentTest,
            {
              state: 'failed',
              title: self.collectorState.getBeforeHookTestTile(),
              isHook: true
            }
          );
        }
      });
    });
  }

  registerAfterAllHooks() {
    const self = this;

    // Logs commands from after all hooks that passed.
    Cypress.mocha.getRunner().on('hook end', function (hook) {
      if (hook.hookName === "after all" && self.collectorState.hasLogsCurrentStack() && !hook._ctr_hook) {
        self.debugLog('extended: sending logs of passed after all hook');
        self.sendLogsToPrinter(
          self.collectorState.getCurrentLogStackIndex(),
          hook,
          {
            state: 'passed',
            title: self.collectorState.getAfterHookTestTile(),
            consoleTitle: self.collectorState.getAfterHookTestTile(),
            isHook: true,
            noQueue: true,
          }
        );
      }
    });

    // Logs after all hook commands when a command fails in the hook.
    Cypress.prependListener('fail', function (error) {
      const currentRunnable = this.mocha.getRunner().currentRunnable;

      if (currentRunnable.hookName === 'after all' && self.collectorState.hasLogsCurrentStack()) {
        // We only have the full list of commands when the suite ends.
        this.mocha.getRunner().prependOnceListener('suite end', () => {
          self.debugLog('extended: sending logs of failed after all hook');
          self.sendLogsToPrinter(
            self.collectorState.getCurrentLogStackIndex(),
            currentRunnable,
            {
              state: 'failed',
              title: self.collectorState.getAfterHookTestTile(),
              isHook: true,
              noQueue: true,
              wait: 5, // Need to wait so that cypress log updates happen.
            }
          );
        });

        // Have to wait for debounce on log updates to have correct state information.
        // Done state is used as callback and awaited in Cypress.fail.
        Cypress.state('done', async (error) => {
          await new Promise(resolve => setTimeout(resolve, 6));
          throw error;
        });
      }

      Cypress.state('error', error);
      throw error;
    });
  }

  registerTests() {
    const self = this;

    const sendLogsToPrinterForATest = (test) => {
      // We take over logging the passing test titles since we need to control when it gets printed so
      // that our logs come after it is printed.
      if (test.state === 'passed') {
        this.printPassingMochaTestTitle(test);
        this.preventNextMochaPassEmit();
      }

      this.sendLogsToPrinter(this.collectorState.getCurrentLogStackIndex(), test, {noQueue: true});
    };

    const testHasAfterEachHooks = (test) => {
      do {
        if (test.parent._afterEach.length > 0) {
          return true;
        }
        test = test.parent;
      } while(test.parent);
      return false;
    };

    const isLastAfterEachHookForTest = (test, hook) => {
      let suite = test.parent;
      do {
        if (suite._afterEach.length === 0) {
          suite = suite.parent;
        } else {
          return suite._afterEach.indexOf(hook) === suite._afterEach.length - 1;
        }
      } while (suite);
      return false;
    };

    // Logs commands form each separate test when after each hooks are present.
    Cypress.mocha.getRunner().on('hook end', function (hook) {
      if (hook.hookName === 'after each') {
        if (isLastAfterEachHookForTest(self.collectorState.getCurrentTest(), hook)) {
          self.debugLog('extended: sending logs for ended test, just after the last after each hook: ' + self.collectorState.getCurrentTest().title);
          sendLogsToPrinterForATest(self.collectorState.getCurrentTest());
        }
      }
    });
    // Logs commands form each separate test when there is no after each hook.
    Cypress.mocha.getRunner().on('test end', function (test) {
      if (!testHasAfterEachHooks(test)) {
        self.debugLog('extended: sending logs for ended test, that has not after each hooks: ' + self.collectorState.getCurrentTest().title);
        sendLogsToPrinterForATest(self.collectorState.getCurrentTest());
      }
    });
    // Logs commands if test was manually skipped.
    Cypress.mocha.getRunner().on('pending', function (test) {
      if (self.collectorState.getCurrentTest() === test) {
        // In case of fully skipped tests we might not yet have a log stack.
        if (self.collectorState.hasLogsCurrentStack()) {
          self.debugLog('extended: sending logs for skipped test: ' + test.title);
          sendLogsToPrinterForATest(test);
        }
      }
    });
  }

  registerLogToFiles() {
    after(function () {
      cy.wait(6, {log: false});
      cy.task(CONSTANTS.TASK_NAME_OUTPUT, null, {log: false});
    });
  }

  debounceNextMochaSuite(promise) {
    const runner = Cypress.mocha.getRunner();

    // Hack to make mocha wait for our logs to be written to console before
    // going to the next suite. This is because 'fail' and 'suite begin' both
    // fire synchronously and thus we wouldn't get a window to display the
    // logs between the failed hook title and next suite title.
    const originalRunSuite = runner.runSuite;
    runner.runSuite = function (...args) {
      promise
        .catch(() => {/* noop */})
        // We need to wait here as for some reason the next suite title will be displayed to soon.
        .then(() => new Promise(resolve => setTimeout(resolve, 6)))
        .then(() => {
          originalRunSuite.apply(runner, args);
          runner.runSuite = originalRunSuite;
        });
    }
  }

  registerCypressBeforeMochaHooksSealEvent() {
    // Hack to have dynamic after hook per suite.
    // The onSpecReady in cypress is called before the hooks are 'condensed', or so
    // to say sealed and thus in this phase we can register dynamically hooks.
    const oldOnSpecReady = Cypress.onSpecReady;
    Cypress.onSpecReady = function () {
      Cypress.emit('before:mocha:hooks:seal');
      oldOnSpecReady(...arguments);
    };
  }

  prependBeforeAllHookInAllSuites(rootSuites, hookCallback) {
    const recursiveSuites = (suites) => {
      if (suites) {
        suites.forEach((suite) => {
          if (suite.isPending()) {
            return
          }
          suite.afterAll(hookCallback);
          // Make sure our hook is first so that other after all hook logs come after
          // the failed before all hooks logs.
          const hook = suite._afterAll.pop();
          suite._afterAll.unshift(hook);
          // Don't count this in the hook index and logs.
          hook._ctr_hook = true;

          recursiveSuites(suite.suites);
        });
      }
    };
    recursiveSuites(rootSuites);
  }

  printPassingMochaTestTitle(test) {
    if (Cypress.config('isTextTerminal')) {
      Cypress.emit('mocha', 'pass', {
        "id": test.id,
        "order": test.order,
        "title": test.title,
        "state": "passed",
        "type": "test",
        "duration": test.duration,
        "wallClockStartedAt": test.wallClockStartedAt,
        "timings": test.timings,
        "file": null,
        "invocationDetails": test.invocationDetails,
        "final": true,
        "currentRetry": test.currentRetry(),
        "retries": test.retries(),
      })
    }
  }

  preventNextMochaPassEmit() {
    const oldAction = Cypress.action;
    Cypress.action = function (actionName, ...args) {
      if (actionName === 'runner:pass') {
        Cypress.action = oldAction;
        return;
      }

      return oldAction.call(Cypress, actionName, ...args);
    };
  }

  debugLog(message) {
    if (this.config.debug) {
      console.log(CONSTANTS.DEBUG_LOG_PREFIX + message);
    }
  }

};
