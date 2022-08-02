var irreducible = require('./irreducible');

module.exports = irreducible('Error', function (x) { return x instanceof Error; });
