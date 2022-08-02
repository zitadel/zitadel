'use strict';
const test = require('ava');
const del = require('del');
const path = require('path');

const fs = require('fs');
const fsu = require('./index.js');

test.after.always(() => del([ './test' ]));

test.serial.cb('write unique file with callbacks', t => {
  fsu.writeFileUnique(path.join('test', 'test{_file###}.txt'), 'test', { force: true }, (err, path) => {
    if (err) {
      t.fail(err);
    } else {
      t.true(path.endsWith('test.txt'));
    }
    t.end();
  });
});

test.serial.cb('write unique file and stream with callbacks', t => {
  const stream = fsu.createWriteStreamUnique(path.join('test', 'test{_stream###}.txt'));
  fs.createReadStream('readme.md').pipe(stream).on('finish', () => {
    t.true(stream.path.endsWith('test_stream001.txt'));
    t.end();
  });
});

