import { Component, Inject } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { Action, FlowType, TriggerType } from 'src/app/proto/generated/zitadel/action_pb';
import { SetTriggerActionsRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-add-flow-dialog',
  templateUrl: './add-flow-dialog.component.html',
  styleUrls: ['./add-flow-dialog.component.scss'],
})
export class AddFlowDialogComponent {
  public flowType?: FlowType.AsObject;
  public actions: Action.AsObject[] = [];
  public typesForSelection: FlowType.AsObject[] = [];
  public triggerTypesForSelection: TriggerType.AsObject[] = [];

  public form!: UntypedFormGroup;
  constructor(
    private toast: ToastService,
    private mgmtService: ManagementService,
    private fb: UntypedFormBuilder,
    public dialogRef: MatDialogRef<AddFlowDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    if (data.flowType) {
      this.flowType = data.flowType as FlowType.AsObject;
      this.getTriggerTypes((data.flowType as FlowType.AsObject).id);
    }

    this.form = this.fb.group({
      triggerType: [data.triggerType ? data.triggerType : '', [requiredValidator]],
      actionIdsList: [data.actions ? (data.actions as Action.AsObject[]).map((a) => a.id) : [], [requiredValidator]],
    });

    this.getActionIds();
  }

  private getTriggerTypes(id: string): Promise<void> {
    return this.mgmtService
      .listFlowTriggerTypes(id)
      .then((resp) => {
        this.triggerTypesForSelection = resp.resultList;
      })
      .catch((error: any) => {
        this.toast.showError(error);
      });
  }

  private getActionIds(): Promise<void> {
    return this.mgmtService
      .listActions()
      .then((resp) => {
        this.actions = resp.resultList;
      })
      .catch((error: any) => {
        this.toast.showError(error);
      });
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    if (this.flowType) {
      const req = new SetTriggerActionsRequest();
      req.setActionIdsList(this.form.get('actionIdsList')?.value);
      req.setFlowType(this.flowType.id);
      req.setTriggerType((this.form.get('triggerType')?.value as TriggerType.AsObject).id);

      this.dialogRef.close(req);
    }
  }
}
