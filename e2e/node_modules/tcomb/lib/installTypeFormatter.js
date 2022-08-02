var isType = require('./isType');
var getTypeName = require('./getTypeName');

var listStyle = {style: 'list-style-type: none; padding: 5px 10px; margin: 0;'};
var listItemStyle = {style: 'padding: 3px 0;'};
var typeStyle = {style: 'font-weight: bolder;'};
var propStyle = {style: 'color: #881391;'};
var metaStyle = {style: 'color: #777;'};

function getKindName(kind) {
  return kind === 'subtype' ? 'refinement' : kind;
}

var TypeFormatter = {
  header: function (x) {
    if (!isType(x)) {
      return null;
    }
    return ['span',
      ['span', typeStyle, getTypeName(x)],
      ' (' + getKindName(x.meta.kind) + ')'
    ];
  },
  hasBody: function (x) {
    return x.meta.kind !== 'irreducible';
  },
  body: function (x) {
    if (x.meta.kind === 'struct' || x.meta.kind === 'interface') {
      var props = Object.keys(x.meta.props).map(function (prop) {
        return ['li', listItemStyle,
          ['span', propStyle, prop + ': '],
          ['object', { object: x.meta.props[prop] }]
        ];
      });
      return ['ol', listStyle].concat(props);
    }
    if (x.meta.kind === 'dict') {
      return ['ol', listStyle,
        ['li', listItemStyle,
          ['span', metaStyle, 'domain: '],
          ['object', { object: x.meta.domain }]
        ],
        ['li', listItemStyle,
          ['span', metaStyle, 'codomain: '],
          ['object', { object: x.meta.codomain }]
        ]
      ];
    }
    if (x.meta.kind === 'list' || x.meta.kind === 'subtype' || x.meta.kind === 'maybe') {
      return ['ol', listStyle,
        ['li', listItemStyle,
          ['span', metaStyle, 'type: '],
          ['object', { object: x.meta.type }]
        ]
      ];
    }
    if (x.meta.kind === 'enums') {
      var enums = Object.keys(x.meta.map).map(function (e) {
        return ['li', listItemStyle,
          ['span', propStyle, e + ': '],
          ['object', { object: x.meta.map[e] }]
        ];
      });
      return ['ol', listStyle].concat(enums);
    }
    if (x.meta.kind === 'union' || x.meta.kind === 'tuple' || x.meta.kind === 'intersection') {
      var types = x.meta.types.map(function (type) {
        return ['li', listItemStyle,
          ['object', { object: type }]
        ];
      });
      return ['ol', listStyle].concat(types);
    }
  }
};

function installTypeFormatter() {
  if (typeof window !== 'undefined') {
    window.devtoolsFormatters = window.devtoolsFormatters || [];
    window.devtoolsFormatters.push(TypeFormatter);
  }
}

installTypeFormatter.TypeFormatter = TypeFormatter;
module.exports = installTypeFormatter;
