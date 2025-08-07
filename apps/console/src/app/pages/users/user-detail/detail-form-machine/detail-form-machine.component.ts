import { Component, DestroyRef, EventEmitter, Input, Output } from '@angular/core';
import { FormBuilder, FormControl } from '@angular/forms';
import { combineLatestWith, distinctUntilChanged, ReplaySubject } from 'rxjs';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { AccessTokenType, MachineUser } from '@zitadel/proto/zitadel/user/v2/user_pb';
import { startWith } from 'rxjs/operators';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'cnsl-detail-form-machine',
  templateUrl: './detail-form-machine.component.html',
  styleUrls: ['./detail-form-machine.component.scss'],
})
export class DetailFormMachineComponent {
  @Input({ required: true }) public set username(username: string) {
    this.username$.next(username);
  }
  @Input({ required: true }) public set user(user: MachineUser) {
    this.user$.next(user);
  }
  @Input() public set disabled(disabled: boolean) {
    this.disabled$.next(disabled);
  }

  private username$ = new ReplaySubject<string>(1);
  private user$ = new ReplaySubject<MachineUser>(1);
  private disabled$ = new ReplaySubject<boolean>(1);

  public machineForm: ReturnType<typeof this.buildForm>;

  @Output() public submitData = new EventEmitter<ReturnType<(typeof this.machineForm)['getRawValue']>>();

  public accessTokenTypes: AccessTokenType[] = [AccessTokenType.BEARER, AccessTokenType.JWT];

  constructor(
    private readonly fb: FormBuilder,
    private readonly destroyRef: DestroyRef,
  ) {
    this.machineForm = this.buildForm();
  }

  private buildForm() {
    const form = this.fb.group({
      username: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      name: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      description: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      accessTokenType: new FormControl(AccessTokenType.BEARER, { nonNullable: true, validators: [requiredValidator] }),
    });

    form.controls.username.disable();
    this.disabled$
      .pipe(startWith(false), distinctUntilChanged(), takeUntilDestroyed(this.destroyRef))
      .subscribe((disabled) => {
        this.toggleFormControl(form.controls.name, disabled);
        this.toggleFormControl(form.controls.description, disabled);
        this.toggleFormControl(form.controls.accessTokenType, disabled);
      });

    this.username$.pipe(combineLatestWith(this.user$), takeUntilDestroyed(this.destroyRef)).subscribe(([username, user]) => {
      this.machineForm.patchValue({ ...user, username });
    });

    return form;
  }

  public toggleFormControl<T>(control: FormControl<T>, disabled: boolean) {
    if (disabled) {
      control.disable();
      return;
    }
    control.enable();
  }

  public submitForm(): void {
    this.submitData.emit(this.machineForm.getRawValue());
  }
}
