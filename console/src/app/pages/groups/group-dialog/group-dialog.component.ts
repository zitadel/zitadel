import { Component, Inject } from '@angular/core';
import { FormBuilder, FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MessageInitShape } from '@bufbuild/protobuf';
import { TranslateModule } from '@ngx-translate/core';
import { Group } from '@zitadel/proto/zitadel/group/v2/group_pb';
import { CreateGroupRequestSchema, UpdateGroupRequestSchema } from '@zitadel/proto/zitadel/group/v2/group_service_pb';
import {
  maxLengthValidator,
  trimmedRequiredValidator,
} from 'src/app/modules/form-field/validators/validators';

const GROUP_NAME_MAX_LENGTH = 200;
const GROUP_DESCRIPTION_MAX_LENGTH = 200;
import { InputModule } from 'src/app/modules/input/input.module';

@Component({
  selector: 'cnsl-group-dialog',
  templateUrl: './group-dialog.component.html',
  styleUrls: ['./group-dialog.component.scss'],
  imports: [MatButtonModule, MatDialogModule, ReactiveFormsModule, TranslateModule, InputModule],
})
export class GroupDialogComponent {
  protected readonly groupForm: ReturnType<typeof this.buildGroupForm>;

  constructor(
    private readonly fb: FormBuilder,
    public dialogRef: MatDialogRef<
      GroupDialogComponent,
      MessageInitShape<typeof CreateGroupRequestSchema | typeof UpdateGroupRequestSchema>
    >,
    @Inject(MAT_DIALOG_DATA) public readonly data: { group?: Group; organizationId: string },
  ) {
    this.groupForm = this.buildGroupForm();

    if (data.group) {
      this.groupForm.patchValue({
        name: data.group.name,
        description: data.group.description,
      });
    }
  }

  private buildGroupForm() {
    return this.fb.group({
      name: new FormControl<string>('', {
        nonNullable: true,
        validators: [trimmedRequiredValidator, maxLengthValidator(GROUP_NAME_MAX_LENGTH)],
      }),
      description: new FormControl<string>('', {
        nonNullable: true,
        validators: [maxLengthValidator(GROUP_DESCRIPTION_MAX_LENGTH)],
      }),
    });
  }

  protected submit(): void {
    if (this.groupForm.invalid) {
      return;
    }
    const { name, description } = this.groupForm.getRawValue();

    if (this.data.group) {
      this.dialogRef.close({
        id: this.data.group.id,
        name,
        description,
      });
      return;
    }

    this.dialogRef.close({
      organizationId: this.data.organizationId,
      name,
      description,
    });
  }
}
