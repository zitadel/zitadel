import { AbstractControl, ValidationErrors, ValidatorFn, Validators } from '@angular/forms';

export function containsSymbolValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /[^a-z0-9]/gi, 'ERRORS.SYMBOLERROR');
}

export function containsNumberValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /[0-9]/g, 'ERRORS.NUMBERERROR');
}

export function containsUpperCaseValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /[A-Z]/g, 'ERRORS.UPPERCASEMISSING');
}

export function containsLowerCaseValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /[a-z]/g, 'ERRORS.LOWERCASEMISSING');
}

export function phoneValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /^($|(\+|00)[0-9 ]+$)/, 'ERRORS.PHONE');
}

export function requiredValidator(c: AbstractControl): ValidationErrors | null {
  return i18nErr(Validators.required(c), 'ERRORS.REQUIRED');
}

export function emailValidator(c: AbstractControl): ValidationErrors | null {
  return i18nErr(Validators.email(c), 'ERRORS.NOTANEMAIL');
}

export function minLengthValidator(minLength: number): ValidatorFn {
  return (c: AbstractControl): ValidationErrors | null => {
    return i18nErr(Validators.minLength(minLength)(c), 'ERRORS.MINLENGTH', { requiredLength: minLength });
  };
}

export function passwordConfirmValidator(passwordControlName: string = 'password') {
  return (c: AbstractControl): ValidationErrors | null => {
    if (!c.parent || !c) {
      return null;
    }
    const pwd = c.parent.get(passwordControlName);
    const cpwd = c;

    if (!pwd || !cpwd) {
      return null;
    }
    if (pwd.value !== cpwd.value) {
      return i18nErr(undefined, 'ERRORS.PWNOTEQUAL');
    }
    return null;
  };
}

function regexpValidator(c: AbstractControl, regexp: RegExp, i18nKey: string): ValidationErrors | null {
  return !c.value || regexp.test(c.value) ? null : i18nErr({ invalid: true }, i18nKey, { regexp: regexp });
}

function i18nErr(err: ValidationErrors | null | undefined, i18nKey: string, params?: any): ValidationErrors | null {
  if (err === null) {
    return null;
  } else {
    return {
      ...err,
      invalid: true,
      [i18nKey.toLowerCase().replaceAll('.', '')]: {
        valid: false,
        i18nKey: i18nKey,
        params: params,
      },
    };
  }
}
