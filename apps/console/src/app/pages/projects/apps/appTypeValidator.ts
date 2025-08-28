import { FormControl } from '@angular/forms';

export function nativeValidator(c: FormControl): any {
  const REGEXP = /([a-zA-Z0-9]*:\/\/)\w+/g;

  if (c.value.startsWith('http://localhost')) {
    return null;
  } else if (c.value.startsWith('https://') || c.value.startsWith('http://')) {
    return {
      invalid: true,
      nativeValidator: {
        valid: false,
      },
    };
  } else if (REGEXP.test(c.value)) {
    return null;
  }
}

export function webValidator(c: FormControl): any {
  if (c.value.startsWith('https://')) {
    return null;
  } else if (c.value.startsWith('http://')) {
    return {
      invalid: false,
      webValidator: {
        valid: true,
        error: 'LOCALHOSTALLOWEDFORTESTING',
      },
    };
  }
}
