const BaseOutputProcessor = require('./BaseOutputProcessor');

const CONSTANTS = require('../constants');
const PADDING = '    ';
const PADDING_LOGS = `${PADDING}`.repeat(6);
const {EOL} = require('os');

module.exports = class TextOutputProcessor extends BaseOutputProcessor {

  constructor(file) {
    super(file);
    this.chunkSeparator = EOL + EOL;
  }

  severityToFont(severity) {
    return {
      [CONSTANTS.SEVERITY.ERROR]: 'X',
      [CONSTANTS.SEVERITY.WARNING]: '!',
      [CONSTANTS.SEVERITY.SUCCESS]: 'K',
    }[severity];
  }

  padTypeText(text) {
    return Array(Math.max(PADDING_LOGS.length - text.length + 1, 0)).join(' ')
      + text;
  }

  write(allMessages) {

    Object.entries(allMessages).forEach(([spec, tests]) => {
      let text = `${spec}:${EOL}`;
      Object.entries(tests).forEach(([test, messages]) => {
        text += `${PADDING}${test}${EOL}`;
        messages.forEach(([type, message, severity]) => {
          text += (this.padTypeText(`${type} (${this.severityToFont(severity)}): `) +
            message.replace(/\n/g, `${EOL}${PADDING_LOGS}`) + EOL).replace(/\s+\n/, '\n');
        });
        text += EOL;
      });

      this.writeSpecChunk(spec, text);
    });
  }

};
