const BaseOutputProcessor = require('./BaseOutputProcessor');

module.exports = class JsonOutputProcessor extends BaseOutputProcessor {

  constructor(file) {
    super(file);
    this.initialContent = "{\n\n}";
    this.chunkSeparator = ',\n';
  }

  write(allMessages) {
    Object.entries(allMessages).forEach(([spec, tests]) => {
      let data = {[spec]: {}};

      Object.entries(tests).forEach(([test, messages]) => {
        data[spec][test] = messages.map(([type, message, severity]) => ({
          type: type,
          severity: severity,
          message: message,
        }))
      });

      let chunk = JSON.stringify(data, null, 2);
      chunk = chunk.slice(2, -2);

      this.writeSpecChunk(spec, chunk, -2);
    });
  }

};
