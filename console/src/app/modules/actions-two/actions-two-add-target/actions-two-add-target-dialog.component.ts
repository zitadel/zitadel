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
import { Target, TargetSchema } from '@zitadel/proto/zitadel/action/v2beta/target_pb';

enum TargetType {
  RestWebhook = 'restWebhook',
  RestCall = 'restCall',
  RestAsync = 'restAsync',
}

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
  public TargetType = TargetType;
  public targetTypeValues = Object.values(TargetType); // Get enum values

  protected readonly targetForm: ReturnType<typeof this.buildTargetForm>;

  constructor(
    private fb: FormBuilder,
    public dialogRef: MatDialogRef<ActionTwoAddTargetDialogComponent, MessageInitShape<typeof TargetSchema>>,
    @Inject(MAT_DIALOG_DATA) public data: { target: Target },
  ) {
    console.log(data.target);

    this.targetForm = this.buildTargetForm();

    if (data.target) {
      this.targetForm.patchValue({
        name: data.target?.name,
        endpoint: data.target.endpoint,
        timeout: Number(data.target.timeout?.seconds),
        interrupt_on_error:
          data.target.targetType.case === 'restWebhook' || data.target.targetType.case === 'restCall'
            ? data.target.targetType.value.interruptOnError
            : false,
      });
    }
  }

  public buildTargetForm() {
    return this.fb.group({
      name: new FormControl<string>('', [requiredValidator]),
      type: new FormControl<TargetType>(TargetType.RestWebhook, [requiredValidator]),
      endpoint: new FormControl<string>('', [requiredValidator]),
      timeout: new FormControl<number>(10, [requiredValidator]),
      interrupt_on_error: new FormControl<boolean>(true),
    });
  }

  public closeWithResult() {
    if (this.targetForm.valid) {
      const timeoutDuration: MessageInitShape<typeof DurationSchema> = {
        seconds: BigInt(this.targetForm.get('timeout')?.value ?? 10),
        nanos: 0,
      };

      let req: MessageInitShape<typeof TargetSchema> = {
        name: this.targetForm.get('name')?.value ?? '',
        endpoint: this.targetForm.get('endpoint')?.value ?? '',
        timeout: timeoutDuration,
        targetType: {
          case: this.targetType as 'restWebhook' | 'restCall',
          value: {
            interruptOnError:
              this.targetType == 'restWebhook' || this.targetType == 'restCall'
                ? !!this.targetForm.get('interrupt_on_error')?.value
                : undefined,
          },
        },
      };

      this.dialogRef.close(req);
    }
  }

  public get targetType(): TargetType {
    return this.targetForm.get('type')?.value as TargetType;
  }
}
