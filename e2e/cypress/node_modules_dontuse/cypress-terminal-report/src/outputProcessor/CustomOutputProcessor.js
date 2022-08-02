const BaseOutputProcessor = require('./BaseOutputProcessor');

module.exports = class CustomOutputProcessor extends BaseOutputProcessor {

  constructor(file, processorCallback) {
    super(file);
    this.processorCallback = processorCallback;
  }

  write(allMessages) {
    this.processorCallback.call(this, allMessages);
  }

};
