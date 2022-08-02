module.exports = {
  toReadableString: function (errorList) {
    return '\n' + errorList.map((error) => {
      return `=> ${error.dataPath.replace(/\//, '.')}: ${error.message}`;
    }).join('\n') + '\n';
  }
};
