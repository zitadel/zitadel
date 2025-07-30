export function _base64ToArrayBuffer(thing: any) {
  if (typeof thing === 'string') {
    // base64url to base64
    thing = thing.replace(/-/g, '+').replace(/_/g, '/');

    // base64 to Uint8Array
    var str = window.atob(thing);
    var bytes = new Uint8Array(str.length);
    for (var i = 0; i < str.length; i++) {
      bytes[i] = str.charCodeAt(i);
    }
    thing = bytes;
  }

  // Array to Uint8Array
  if (Array.isArray(thing)) {
    thing = new Uint8Array(thing);
  }

  // Uint8Array to ArrayBuffer
  if (thing instanceof Uint8Array) {
    thing = thing.buffer;
  }

  // error if none of the above worked
  if (!(thing instanceof ArrayBuffer)) {
    throw new TypeError('could not coerce to ArrayBuffer');
  }

  return thing;
}
