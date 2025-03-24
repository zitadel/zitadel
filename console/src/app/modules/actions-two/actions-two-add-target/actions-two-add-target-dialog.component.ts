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
import { GetTarget } from '@zitadel/proto/zitadel/resources/action/v3alpha/target_pb';

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
  ],
})
export class ActionTwoAddTargetDialogComponent {
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
      endpoint: new FormControl<string>('', [requiredValidator]),
      timeout: new FormControl<number>(10, [requiredValidator]),
      interrupt_on_error: new FormControl<boolean>(true),
      // await_response: new FormControl<boolean>(false),
    });
  }

  public closeWithResult() {
    if (this.targetForm.valid) {
      const timeoutDuration: MessageInitShape<typeof DurationSchema> = {
        seconds: BigInt(this.targetForm.get('timeout')?.value ?? 10),
        nanos: 0,
      };

      let req: MessageInitShape<typeof PatchTargetRequestSchema> | MessageInitShape<typeof CreateTargetRequestSchema>;
      if (this.data.target) {
        req = {
          // instance_id: this.data.instance_id,
          target: {
            name: this.targetForm.get('name')?.value ?? '',
            endpoint: this.targetForm.get('endpoint')?.value ?? '',
            timeout: timeoutDuration,
            targetType: {
              case: 'restWebhook',
              value: {
                interruptOnError: !!this.targetForm.get('interrupt_on_error')?.value,
              },
            },
            // await_response: this.targetForm.value.await_response,
          },
        };
      } else {
        req = {
          // instance_id: this.data.instance_id,
          target: {
            name: this.targetForm.get('name')?.value ?? '',
            endpoint: this.targetForm.get('endpoint')?.value ?? '',
            timeout: timeoutDuration,
            targetType: {
              case: 'restWebhook',
              value: {
                interruptOnError: !!this.targetForm.get('interrupt_on_error')?.value,
              },
            },
            // await_response: this.targetForm.value.await_response,
          },
        };
      }

      this.dialogRef.close(req);
    }
  }
}
