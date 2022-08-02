"use strict";

const t = require('tcomb');

const {
  isUUID,
  isISO8601
} = require('validator');

const TestState = t.enums.of(['passed', 'failed', 'pending', 'skipped'], 'TestState');
const TestSpeed = t.enums.of(['slow', 'medium', 'fast'], 'TestSpeed');
const DateString = t.refinement(t.String, isISO8601, 'DateString');
const Duration = t.maybe(t.Integer);
const Uuid = t.refinement(t.String, isUUID, 'UUID');
const ReportMeta = t.struct({
  mocha: t.struct({
    version: t.String
  }),
  mochawesome: t.struct({
    options: t.Object,
    version: t.String
  }),
  marge: t.struct({
    options: t.Object,
    version: t.String
  })
});
const Test = t.struct({
  title: t.String,
  fullTitle: t.String,
  timedOut: t.maybe(t.Boolean),
  duration: Duration,
  state: t.maybe(TestState),
  speed: t.maybe(TestSpeed),
  pass: t.Boolean,
  fail: t.Boolean,
  pending: t.Boolean,
  code: t.String,
  err: t.Object,
  uuid: Uuid,
  parentUUID: t.maybe(Uuid),
  skipped: t.Boolean,
  context: t.maybe(t.String),
  isHook: t.Boolean
});
const Suite = t.declare('Suite');
Suite.define(t.struct({
  title: t.String,
  suites: t.list(Suite),
  tests: t.list(Test),
  root: t.Boolean,
  _timeout: t.Integer,
  file: t.String,
  uuid: Uuid,
  fullFile: t.String,
  beforeHooks: t.list(Test),
  afterHooks: t.list(Test),
  passes: t.list(Uuid),
  failures: t.list(Uuid),
  pending: t.list(Uuid),
  skipped: t.list(Uuid),
  duration: Duration,
  rootEmpty: t.maybe(t.Boolean)
}));
const TestReport = t.struct({
  stats: t.struct({
    suites: t.Integer,
    tests: t.Integer,
    passes: t.Integer,
    pending: t.Integer,
    failures: t.Integer,
    start: DateString,
    end: DateString,
    duration: Duration,
    testsRegistered: t.Integer,
    passPercent: t.Number,
    pendingPercent: t.Number,
    other: t.Integer,
    hasOther: t.Boolean,
    skipped: t.Integer,
    hasSkipped: t.Boolean
  }),
  results: t.list(Suite),
  meta: t.maybe(ReportMeta)
});
module.exports = {
  TestReport,
  Test,
  Suite
};