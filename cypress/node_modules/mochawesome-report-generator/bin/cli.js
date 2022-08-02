#!/usr/bin/env node

/* eslint-disable no-unused-expressions */
"use strict";

const yargs = require('yargs');

const {
  yargsOptions
} = require('../lib/options');

const marge = require('./cli-main'); // Setup yargs


yargs.usage('Usage: $0 [options] <data_file> [data_file2..]').demand(1).options(yargsOptions).help('help').alias('h', 'help').version().epilog(`Copyright ${new Date().getFullYear()} Adam Gruber`).argv; // Call the main cli program

marge(yargs.argv);