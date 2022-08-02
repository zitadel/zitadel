module.exports = class LogFormat {

  constructor(config) {
    this.config = config;
  }

  formatXhrLog(xhrLog) {
    let logMessage = '';
    if (xhrLog.response) {
      logMessage += `Status: ${xhrLog.response.status}\n`;
    } else if (xhrLog.networkError) {
      logMessage += `Network error: ${xhrLog.networkError}\n`;
    }

    if (xhrLog.request) {
      if (this.config.collectRequestData) {
        if (this.config.collectHeaderData) {
          logMessage += `Request headers: ${xhrLog.request.headers}\n`;
        }
        logMessage += `Request body: ${xhrLog.request.body}\n`;
      }
    }

    if (xhrLog.response) {
      if (this.config.collectHeaderData) {
        logMessage += `Response headers: ${xhrLog.response.headers}\n`;
      }
      logMessage += `Response body: ${xhrLog.response.body}`;
    }

    return logMessage.trimEnd();
  }

  formatXhrBody(body) {
    if (!body) {
      return Promise.resolve('<EMPTY>');
    } else if (typeof body === 'string') {
      return Promise.resolve(body);
    } else if (typeof body === 'object') {
      if (typeof body.text === 'function') {
        return body.text();
      }
      return Promise.resolve(`${JSON.stringify(body, null, 2)}`);
    }
    return Promise.resolve('<UNKNOWN>');
  }

}
