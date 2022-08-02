module.exports = {
  TASK_NAME: 'ctrLogMessages',
  TASK_NAME_OUTPUT: 'ctrLogFiles',

  LOG_TYPES: {
    BROWSER_CONSOLE_LOG: 'cons:log',
    BROWSER_CONSOLE_INFO: 'cons:info',
    BROWSER_CONSOLE_WARN: 'cons:warn',
    BROWSER_CONSOLE_ERROR: 'cons:error',
    BROWSER_CONSOLE_DEBUG: 'cons:debug',

    CYPRESS_LOG: 'cy:log',
    CYPRESS_XHR: 'cy:xhr',
    CYPRESS_FETCH: 'cy:fetch',
    CYPRESS_REQUEST: 'cy:request',
    CYPRESS_ROUTE: 'cy:route',
    CYPRESS_INTERCEPT: 'cy:intercept',
    CYPRESS_COMMAND: 'cy:command',

    PLUGIN_LOG_TYPE: 'ctr:info',
  },

  SEVERITY: {
    SUCCESS: 'success',
    ERROR: 'error',
    WARNING: 'warning',
  },

  HOOK_TITLES: {
    BEFORE: '[[ before all {index} ]]',
    AFTER: '[[ after all {index} ]]',
  },

  PADDING: {
    LOG: Array(21).join(' '),
  },

  DEBUG_LOG_PREFIX: 'CTR-DEBUG: ',

  // HTTP methods defined by the `node:http` module
  HTTP_METHODS: [
    'ACL',         'BIND',       'CHECKOUT',
    'CONNECT',     'COPY',       'DELETE',
    'GET',         'HEAD',       'LINK',
    'LOCK',        'M-SEARCH',   'MERGE',
    'MKACTIVITY',  'MKCALENDAR', 'MKCOL',
    'MOVE',        'NOTIFY',     'OPTIONS',
    'PATCH',       'POST',       'PROPFIND',
    'PROPPATCH',   'PURGE',      'PUT',
    'REBIND',      'REPORT',     'SEARCH',
    'SOURCE',      'SUBSCRIBE',  'TRACE',
    'UNBIND',      'UNLINK',     'UNLOCK',
    'UNSUBSCRIBE'
  ]
};
