import { Injectable } from '@angular/core';
import { ValidatorFn } from '@angular/forms';
import {
  containsLowerCaseValidator,
  containsNumberValidator,
  containsSymbolValidator,
  containsUpperCaseValidator,
  minLengthValidator,
  requiredValidator,
} from '../modules/form-field/validators/validators';
import { PasswordComplexityPolicy } from '@zitadel/proto/zitadel/policy_pb';

@Injectable({
  providedIn: 'root',
})
export class PasswordComplexityValidatorFactoryService {
  constructor() {}

  buildValidators(policy?: PasswordComplexityPolicy) {
    const validators: ValidatorFn[] = [requiredValidator];
    if (policy?.minLength) {
      validators.push(minLengthValidator(Number(policy.minLength)));
    }
    if (policy?.hasLowercase) {
      validators.push(containsLowerCaseValidator);
    }
    if (policy?.hasUppercase) {
      validators.push(containsUpperCaseValidator);
    }
    if (policy?.hasNumber) {
      validators.push(containsNumberValidator);
    }
    if (policy?.hasSymbol) {
      validators.push(containsSymbolValidator);
    }
    return validators;
  }
}
