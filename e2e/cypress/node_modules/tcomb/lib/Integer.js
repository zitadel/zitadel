var refinement = require('./refinement');
var Number = require('./Number');

module.exports = refinement(Number, function (x) { return x % 1 === 0; }, 'Integer');
