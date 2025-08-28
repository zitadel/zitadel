import { Component, Inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { InputModule } from '../../input/input.module';
import { requiredValidator } from '../../form-field/validators/validators';
import { MessageInitShape } from '@bufbuild/protobuf';
import { DurationSchema } from '@bufbuild/protobuf/wkt';
import { MatSelectModule } from '@angular/material/select';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';
import {
  CreateTargetRequestSchema,
  UpdateTargetRequestSchema,
} from '@zitadel/proto/zitadel/action/v2beta/action_service_pb';

type TargetTypes = ActionTwoAddTargetDialogComponent['targetTypes'][number];

@Component({
  selector: 'cnsl-actions-two-add-target-dialog',
  templateUrl: './actions-two-add-target-dialog.component.html',
  styleUrls: ['./actions-two-add-target-dialog.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatDialogModule,
    ReactiveFormsModule,
    TranslateModule,
    InputModule,
    MatCheckboxModule,
    MatSelectModule,
  ],
})
export class ActionTwoAddTargetDialogComponent {
  protected readonly targetTypes = ['restCall', 'restWebhook', 'restAsync'] as const;
  protected readonly targetForm: ReturnType<typeof this.buildTargetForm>;

  constructor(
    private fb: FormBuilder,
    public dialogRef: MatDialogRef<
      ActionTwoAddTargetDialogComponent,
      MessageInitShape<typeof CreateTargetRequestSchema | typeof UpdateTargetRequestSchema>
    >,
    @Inject(MAT_DIALOG_DATA) public readonly data: { target?: Target },
  ) {
    this.targetForm = this.buildTargetForm();

    if (!data?.target) {
      return;
    }

    this.targetForm.patchValue({
      name: data.target.name,
      endpoint: data.target.endpoint,
      timeout: Number(data.target.timeout?.seconds),
      type: this.data.target?.targetType?.case ?? 'restWebhook',
      interruptOnError:
        data.target.targetType.case === 'restWebhook' || data.target.targetType.case === 'restCall'
          ? data.target.targetType.value.interruptOnError
          : false,
    });
  }

  public buildTargetForm() {
    return this.fb.group({
      name: new FormControl<string>('', { nonNullable: true, validators: [requiredValidator] }),
      type: new FormControl<TargetTypes>('restWebhook', {
        nonNullable: true,
        validators: [requiredValidator],
      }),
      endpoint: new FormControl<string>('', { nonNullable: true, validators: [requiredValidator] }),
      timeout: new FormControl<number>(10, { nonNullable: true, validators: [requiredValidator] }),
      interruptOnError: new FormControl<boolean>(false, { nonNullable: true }),
    });
  }

  public closeWithResult() {
    if (this.targetForm.invalid) {
      return;
    }

    const { type, name, endpoint, timeout, interruptOnError } = this.targetForm.getRawValue();

    const timeoutDuration: MessageInitShape<typeof DurationSchema> = {
      seconds: BigInt(timeout),
      nanos: 0,
    };

    const targetType: MessageInitShape<typeof CreateTargetRequestSchema>['targetType'] =
      type === 'restWebhook'
        ? { case: type, value: { interruptOnError } }
        : type === 'restCall'
          ? { case: type, value: { interruptOnError } }
          : { case: 'restAsync', value: {} };

    const baseReq = {
      name,
      endpoint,
      timeout: timeoutDuration,
      targetType,
    };

    this.dialogRef.close(
      this.data.target
        ? {
            ...baseReq,
            id: this.data.target.id,
          }
        : baseReq,
    );
  }
}
