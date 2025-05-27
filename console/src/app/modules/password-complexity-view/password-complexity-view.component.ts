import { Component, Input } from '@angular/core';
import { AbstractControl } from '@angular/forms';
import { PasswordComplexityPolicy } from '@zitadel/proto/zitadel/policy_pb';

@Component({
  selector: 'cnsl-password-complexity-view',
  templateUrl: './password-complexity-view.component.html',
  styleUrls: ['./password-complexity-view.component.scss'],
})
export class PasswordComplexityViewComponent {
  @Input() public password: AbstractControl | null = null;
  @Input({ required: true }) public policy!: PasswordComplexityPolicy;

  protected get minLength() {
    return Number(this.policy.minLength);
  }
}
