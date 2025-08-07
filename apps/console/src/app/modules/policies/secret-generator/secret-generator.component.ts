import { Component } from '@angular/core';
import { SecretGenerator, SecretGeneratorType } from 'src/app/proto/generated/zitadel/settings_pb';

@Component({
  selector: 'cnsl-secret-generator',
  templateUrl: './secret-generator.component.html',
  styleUrls: ['./secret-generator.component.scss'],
})
export class SecretGeneratorComponent {
  public generators: SecretGenerator.AsObject[] = [];

  public readonly AVAILABLEGENERATORS: SecretGeneratorType[] = [
    SecretGeneratorType.SECRET_GENERATOR_TYPE_INIT_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_OTP_SMS,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_OTP_EMAIL,
  ];

  constructor() {}
}
