import { Component, Input } from '@angular/core';
import { FormControl } from '@angular/forms';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';

@Component({
  selector: 'cnsl-password-complexity-view',
  templateUrl: './password-complexity-view.component.html',
  styleUrls: ['./password-complexity-view.component.scss'],
})
export class PasswordComplexityViewComponent {
  @Input() public password!: FormControl;
  @Input() public policy!: PasswordComplexityPolicy.AsObject;
}
