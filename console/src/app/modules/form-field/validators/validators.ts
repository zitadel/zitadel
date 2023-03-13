import { AbstractControl, UntypedFormControl, ValidationErrors, Validators } from '@angular/forms';

export function symbolValidator(c: UntypedFormControl): any {
  const REGEXP: RegExp = /[^a-z0-9]/gi;

  return REGEXP.test(c.value)
    ? null
    : {
        invalid: true,
        symbolValidator: {
          valid: false,
        },
      };
}

export function numberValidator(c: UntypedFormControl): any {
  const REGEXP = /[0-9]/g;

  return REGEXP.test(c.value)
    ? null
    : {
        invalid: true,
        numberValidator: {
          valid: false,
        },
      };
}

export function upperCaseValidator(c: UntypedFormControl): any {
  const REGEXP = /[A-Z]/g;

  return REGEXP.test(c.value)
    ? null
    : {
        invalid: true,
        upperCaseValidator: {
          valid: false,
        },
      };
}

export function lowerCaseValidator(c: UntypedFormControl): any {
  const REGEXP = /[a-z]/g;

  return REGEXP.test(c.value)
    ? null
    : {
        invalid: true,
        lowerCaseValidator: {
          valid: false,
        },
      };
}

export function phoneValidator(c: AbstractControl): ValidationErrors | null {
  const REGEXP = /^($|(\+|00)[0-9 ]+$)/;

  return !c.value || REGEXP.test(c.value)
    ? null
    : {
        invalid: true,
        phoneValidator: {
          valid: false,
          i18nKey: "ERRORS.INVALID_FORMAT"
        },
      };
}

export function requiredValidator(c: AbstractControl): ValidationErrors | null {
  let err = Validators.required(c)
  if (err) {
    err = {
      ...err,
      invalid: true,
      requiredValidator: {
        valid: false,
        i18nKey: "ERRORS.REQUIRED"
      }
    }
  }
  return err
}