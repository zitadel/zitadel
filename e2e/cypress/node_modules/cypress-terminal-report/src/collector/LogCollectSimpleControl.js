const CONSTANTS = require('../constants');
const LogCollectBaseControl = require('./LogCollectBaseControl');

/**
 * Collects and dispatches all logs from all tests and hooks.
 */
module.exports = class LogCollectSimpleControl extends LogCollectBaseControl {

  constructor(collectorState, config) {
    super();
    this.config = config;
    this.collectorState = collectorState;
  }

  register() {
    this.registerState();
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
      Promise.resolve().then(() => {
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
            continuous: this.config.enableContinuousLogging,
          }
        })
          // For some reason cypress throws empty error although the task indeed works.
          .catch((error) => {/* noop */})
      }).catch(console.error);
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
              continuous: this.config.enableContinuousLogging,
            },
            {log: false}
          );
        });
    }
  }

  registerState() {
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
  }

  registerTests() {
    const self = this;

    if (this.config.enableContinuousLogging) {
      this.collectorState.on('log', () => {
        self.sendLogsToPrinter(self.collectorState.getCurrentLogStackIndex(), self.collectorState.getCurrentTest(), {noQueue: true});
        this.collectorState.addNewLogStack();
      });
      return;
    }

    afterEach(function () {
      self.sendLogsToPrinter(self.collectorState.getCurrentLogStackIndex(), self.collectorState.getCurrentTest());
    });

    // Logs commands if test was manually skipped.
    Cypress.mocha.getRunner().on('pending', function (test) {
      if (self.collectorState.getCurrentTest()) {
        // In case of fully skipped tests we might not yet have a log stack.
        if (!self.collectorState.hasLogsCurrentStack()) {
          self.collectorState.addNewLogStack();
        }
        self.sendLogsToPrinter(self.collectorState.getCurrentLogStackIndex(), self.collectorState.getCurrentTest(), {noQueue: true});
      }
    });
  }

  registerLogToFiles() {
    after(function () {
      // Need to wait otherwise some last commands get omitted from logs.
      cy.task(CONSTANTS.TASK_NAME_OUTPUT, null, {log: false});
    });
  }

};
