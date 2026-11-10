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

export function atLeastOneFieldValidator(fields: string[]): ValidatorFn {
  return (formGroup: AbstractControl): ValidationErrors | null => {
    const isValid = fields.some((field) => {
      const control = formGroup.get(field);
      return control && control.value;
    });

    return isValid ? null : { atLeastOneRequired: true }; // Return an error if none are set
  };
}

export function minArrayLengthValidator(minArrLength: number): ValidatorFn {
  return (c: AbstractControl): ValidationErrors | null => {
    return arrayLengthValidator(c, minArrLength, 'ERRORS.ATLEASTONE');
  };
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

function arrayLengthValidator(c: AbstractControl, length: number, i18nKey: string): ValidationErrors | null {
  const arr: string[] = c.value;
  const invalidStrings: string[] = arr.filter((val: string) => val.trim() === '');
  return arr && invalidStrings.length === 0 && arr.length >= length ? null : i18nErr({ invalid: true }, i18nKey);
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

const isFieldEmpty = (fieldName: string, g: AbstractControl) => {
  const field = g.get(fieldName)?.value;
  if (typeof field === 'number') {
    return field && field >= 0 ? true : false;
  }
  if (typeof field === 'string') {
    return field && field.length > 0 ? true : false;
  }
  return false;
};

// Reference: https://stackoverflow.com/a/56057955
export function atLeastOneIsFilled(...fields: string[]): ValidationErrors | null {
  return (g: AbstractControl): ValidationErrors | null => {
    return fields.some((fieldName) => isFieldEmpty(fieldName, g))
      ? null
      : ({ atLeastOne: 'At least one field has to be provided.' } as ValidationErrors);
  };
}
