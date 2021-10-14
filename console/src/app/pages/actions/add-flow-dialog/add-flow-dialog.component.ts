import { Component, Inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
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
    public actions: Action.AsObject[] = [];
    public typesForSelection: FlowType[] = [
      FlowType.FLOW_TYPE_EXTERNAL_AUTHENTICATION,
    ];
    public triggerTypesForSelection: TriggerType[] = [
      TriggerType.TRIGGER_TYPE_POST_AUTHENTICATION,
      TriggerType.TRIGGER_TYPE_POST_CREATION,
      TriggerType.TRIGGER_TYPE_PRE_CREATION,      
    ];
    
    public form!: FormGroup;
    constructor(
      private toast: ToastService,
      private mgmtService: ManagementService,
      private fb: FormBuilder,
      public dialogRef: MatDialogRef<AddFlowDialogComponent>,
      @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
      if (data && data.actionIds) {
        this.actions = data.actionIds;
      }

      this.form = this.fb.group({
        flowType: [data.flowType ? data.flowType : '', [Validators.required]],
        triggerType: [data.triggerType ? data.triggerType : '', [Validators.required]],
        actionIdsList: [this.actions, [Validators.required]],
      });

      this.getActionIds();
    }

    private getActionIds(): void {
      this.mgmtService.listActions().then(resp => {
        this.actions = resp.resultList;
      }).catch((error: any) => {
        this.toast.showError(error);
      });
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
      // if (this.id) {
        // const req = new UpdateActionRequest();
        // req.setId(this.id);
        // req.setName(this.name);
        // req.setScript(this.script);

        // const duration = new Duration();
        // duration.setNanos(0);
        // duration.setSeconds(this.durationInSec);

        // req.setAllowedToFail(this.allowedToFail);

        // req.setTimeout(duration)
        // this.dialogRef.close(req);
      // } else {
        const req = new SetTriggerActionsRequest();
        req.setActionIdsList(this.form.get('actionIdsList')?.value);
        req.setFlowType(this.form.get('flowType')?.value);
        req.setTriggerType(this.form.get('triggerType')?.value);

        this.dialogRef.close(req);
      // }
    }
}
