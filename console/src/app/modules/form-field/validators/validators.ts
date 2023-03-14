import { AbstractControl, ValidationErrors, ValidatorFn, Validators } from '@angular/forms';

export function containsSymbolValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /[^a-z0-9]/gi, "ERRORS.CONTAINSSYMBOL")
}

export function containsNumberValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /[0-9]/g, "ERRORS.CONTAINSNUMBER")
}

export function containsUpperCaseValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /[A-Z]/g, "ERRORS.CONTAINSUPPERCASE")
}

export function containsLowerCaseValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /[a-z]/g, "ERRORS.CONTAINSLOWERCASE")
}

export function phoneValidator(c: AbstractControl): ValidationErrors | null {
  return regexpValidator(c, /^($|(\+|00)[0-9 ]+$)/, "ERRORS.PHONE")
}

export function requiredValidator(c: AbstractControl): ValidationErrors | null {
  return i18nErr(Validators.required(c), 'ERRORS.REQUIRED');
}


export function emailValidator(c: AbstractControl): ValidationErrors | null {
  return i18nErr(Validators.email(c), "ERRORS.EMAIL");
}

export function minLengthValidator(minLength: number): ValidatorFn {
  return (c: AbstractControl): ValidationErrors | null  => {
    return i18nErr(Validators.minLength(minLength)(c), 'ERRORS.MINLENGTH');
  }
}

export function passwordConfirmValidator(passwordControlName: string){
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
      return i18nErr(null, 'USER.PASSWORD.NOTEQUAL')
    }
    return null
  }
}

function regexpValidator(c: AbstractControl, regexp: RegExp, i18nKey: string): ValidationErrors | null {
  return !c.value || regexp.test(c.value)
  ? null
  : i18nErr({invalid: true}, i18nKey)
}

function i18nErr(err: ValidationErrors | null, i18nKey: string): ValidationErrors | null{
  if (err) {
    err = {
      ...err,
      invalid: true,
      [i18nKey.toLowerCase().replaceAll(".","")]: {
        valid: false,
        i18nKey: i18nKey,
      },
    };
  }
  return err;
}
