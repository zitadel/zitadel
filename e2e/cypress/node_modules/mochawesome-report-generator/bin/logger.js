"use strict";

/**
 * Logger - wraps console methods
 *
 */
module.exports = {
  info: msg => console.log.apply(console, [msg]),
  error: msg => console.error.apply(console, [msg])
};