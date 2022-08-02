var irreducible = require('./irreducible');

module.exports = irreducible('Date', function (x) { return x instanceof Date; });
