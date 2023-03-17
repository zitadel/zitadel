import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { Subscription } from 'rxjs';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { AccessTokenType, Human, Machine } from 'src/app/proto/generated/zitadel/user_pb';

@Component({
  selector: 'cnsl-detail-form-machine',
  templateUrl: './detail-form-machine.component.html',
  styleUrls: ['./detail-form-machine.component.scss'],
})
export class DetailFormMachineComponent implements OnInit, OnDestroy {
  @Input() public username!: string;
  @Input() public user!: Human.AsObject | Machine.AsObject;
  @Input() public disabled: boolean = false;
  @Output() public submitData: EventEmitter<any> = new EventEmitter<any>();

  public machineForm!: UntypedFormGroup;

  public accessTokenTypes: AccessTokenType[] = [
    AccessTokenType.ACCESS_TOKEN_TYPE_BEARER,
    AccessTokenType.ACCESS_TOKEN_TYPE_JWT,
  ];

  private sub: Subscription = new Subscription();

  constructor(private fb: UntypedFormBuilder) {
    this.machineForm = this.fb.group({
      userName: [{ value: '', disabled: true }, [requiredValidator]],
      name: [{ value: '', disabled: this.disabled }, requiredValidator],
      description: [{ value: '', disabled: this.disabled }],
      accessTokenType: [AccessTokenType.ACCESS_TOKEN_TYPE_BEARER, [requiredValidator]],
    });
  }

  public ngOnInit(): void {
    this.machineForm.patchValue({ ...this.user, userName: this.username });
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public submitForm(): void {
    this.submitData.emit(this.machineForm.value);
  }

  public get name(): AbstractControl | null {
    return this.machineForm.get('name');
  }

  public get userName(): AbstractControl | null {
    return this.machineForm.get('userName');
  }
}
