import { AfterViewInit, Component, Inject, signal, ViewChild } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogClose,
  MatDialogContent,
  MatDialogModule,
  MatDialogRef,
  MatDialogTitle,
} from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import {
  CreateTargetRequestSchema,
  PatchTargetRequestSchema,
  SetExecutionRequestSchema,
} from '@zitadel/proto/zitadel/resources/action/v3alpha/action_service_pb';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { InputModule } from '../../input/input.module';
import { requiredValidator } from '../../form-field/validators/validators';
import { MessageInitShape } from '@bufbuild/protobuf';
import { DurationSchema } from '@bufbuild/protobuf/wkt';
import { GetTarget, Target } from '@zitadel/proto/zitadel/resources/action/v3alpha/target_pb';
import { MatSelectModule } from '@angular/material/select';

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
    public dialogRef: MatDialogRef<
      ActionTwoAddTargetDialogComponent,
      MessageInitShape<typeof CreateTargetRequestSchema> | MessageInitShape<typeof PatchTargetRequestSchema>
    >,
    @Inject(MAT_DIALOG_DATA) public data: { target: GetTarget },
  ) {
    console.log(data.target);

    this.targetForm = this.buildTargetForm();

    if (data.target) {
      this.targetForm.patchValue({
        name: data.target.config?.name,
        endpoint: data.target.config?.endpoint,
        timeout: Number(data.target.config?.timeout?.seconds),
        interrupt_on_error:
          data.target.config?.targetType.case === 'restWebhook' || data.target.config?.targetType.case === 'restCall'
            ? data.target.config?.targetType.value.interruptOnError
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

      let req: MessageInitShape<typeof PatchTargetRequestSchema> | MessageInitShape<typeof CreateTargetRequestSchema> = {
        target: {
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
        },
      };

      this.dialogRef.close(req);
    }
  }

  public get targetType(): TargetType {
    return this.targetForm.get('type')?.value as TargetType;
  }
}
