'use strict';
const fs = require('fs');
const path = require('path');
const inherits = require('util').inherits;
const WriteStream = fs.WriteStream;

const rxFileParts = /(.*)\{([^#{}]*)(#+)([^#{}]*)\}(.*)/;

const defaultDirMode = parseInt('0777', 8) & (~process.umask());
const defaultFileMode = parseInt('0666', 8) & (~process.umask());

const padNum = function(n, width, z) {
  z = z || '0';
  n += '';
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
};

const writeAll = function(fd, buffer, offset, length, position, cb) {
  fs.write(fd, buffer, offset, length, position, (writeErr, written) => {
    if (writeErr) {
      fs.close(fd, () => cb(writeErr));
    } else if (written === length) {
      fs.close(fd, cb);
    } else {
      offset += written;
      length -= written;
      position += written;
      writeAll(fd, buffer, offset, length, position, cb);
    }
  });
};

const mkdirp = function(p, mode, cb) {
  fs.mkdir(p, mode, err => {
    if (!err) {
      cb();
    } else if (err.code === 'ENOENT') {
      mkdirp(path.dirname(p), mode, err => {
        if (err) {
          cb(err);
        } else {
          mkdirp(p, mode, cb);
        }
      });
    } else if (err.code === 'EEXIST') {
      cb();
    } else {
      cb(err);
    }
  });
};

const openUniqueHandler = function(tryNum, fileParts, options, cb) {
  const file = options.simple ? fileParts.tail : tryNum ? (fileParts.head + fileParts.padLeft + padNum(tryNum, fileParts.pad) + fileParts.padRight + fileParts.tail) : (fileParts.head + fileParts.tail);
  const newPath = path.join(fileParts.path, file);

  fs.open(newPath, options.flags || 'w', options.mode || defaultFileMode, (err, fd) => {
    if (err && err.code === 'EEXIST' && !options.simple) {
      openUniqueHandler(++tryNum, fileParts, options, cb);
    } else if (err && err.code === 'ENOENT' && options.force) {
      mkdirp(fileParts.path, defaultDirMode, ere => {
        if (ere) {
          cb(ere);
        } else {
          openUniqueHandler(tryNum, fileParts, options, cb);
        }
      });
    } else {
      cb(err, fd, newPath);
    }
  });
};

const openUnique = function(file, options, cb) {
  file = path.resolve(file);
  const filePath = path.dirname(file);
  const fileName = path.basename(file);

  const fileParts = rxFileParts.exec(fileName);

  if (!fileParts) {
    options.simple = true;
    openUniqueHandler(0, {
      path: filePath,
      tail: fileName
    }, options, cb);
  } else {
    options.simple = false;
    options.flags = 'wx';
    openUniqueHandler(0, {
      path: filePath,
      head: fileParts[1] || '',
      padLeft: fileParts[2],
      pad: fileParts[3].length,
      padRight: fileParts[4],
      tail: fileParts[5] || ''
    }, options, cb);
  }
};

const writeFileUnique = function(filename, data, options, cb) {
  if (cb === undefined) {
    cb = options;
    options = {
      encoding: 'utf8',
      mode: defaultFileMode,
      flags: 'w'
    };
  }

  openUnique(filename, options, (err, fd, newPath) => {
    if (err) {
      cb(err);
    } else {
      const buffer = Buffer.isBuffer(data) ? data : Buffer.from(String(data), options.encoding || 'utf8');
      writeAll(fd, buffer, 0, buffer.length, 0, err => cb(err, newPath));
    }
  });
};

// stream
const WriteStreamUnique = function(file, options) {
  if (options && options.force) {
    this.force = options.force;
    delete options.force;
  }
  WriteStream.call(this, file, options);
};
inherits(WriteStreamUnique, WriteStream);

WriteStreamUnique.prototype.open = function() {
  openUnique(this.path, {
    flags: this.flags,
    mode: this.mode,
    force: this.force
  }, (err, fd, newPath) => {
    if (err) {
      this.destroy();
      this.emit('error', err);
      return;
    }
    this.path = newPath;
    this.fd = fd;
    this.emit('open', fd);
  });
};

const createWriteStreamUnique = function(file, options) {
  return new WriteStreamUnique(file, options);
};

module.exports = {
  openUnique,
  writeFileUnique,
  WriteStreamUnique,
  createWriteStreamUnique
};
