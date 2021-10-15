import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Action, Flow, FlowType, TriggerType } from 'src/app/proto/generated/zitadel/action_pb';
import { SetTriggerActionsRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddFlowDialogComponent } from './add-flow-dialog/add-flow-dialog.component';

@Component({
  selector: 'cnsl-actions',
  templateUrl: './actions.component.html',
  styleUrls: ['./actions.component.scss']
})
export class ActionsComponent implements OnInit {
  public flow!: Flow.AsObject;
  public flowType: FlowType =  FlowType.FLOW_TYPE_EXTERNAL_AUTHENTICATION;

  public typeControl: FormControl = new FormControl(FlowType.FLOW_TYPE_EXTERNAL_AUTHENTICATION);

  public typesForSelection: FlowType[] = [
    FlowType.FLOW_TYPE_EXTERNAL_AUTHENTICATION,
  ];

  public selection: Action.AsObject[] = [];
  
  constructor(
    private mgmtService: ManagementService,
    private dialog: MatDialog,
    private toast: ToastService,
    ) { 
    this.loadFlow();
  }

  private loadFlow() {
    this.mgmtService.getFlow(this.flowType).then(flowResponse => {
      if (flowResponse.flow)
      this.flow = flowResponse.flow;
      console.log(this.flow);
    })
  }

  ngOnInit(): void {
  }

  public clearFlow(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.CLEAR',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'FLOWS.DIALOG.CLEAR.TITLE',
          descriptionKey: 'FLOWS.DIALOG.CLEAR.DESCRIPTION',
        },
        width: '400px',
      });

      dialogRef.afterClosed().subscribe(resp => {
          if (resp) {
        this.mgmtService.clearFlow(this.flowType).then(resp => {
          this.loadFlow();
        }).catch((error: any) => {
          this.toast.showError(error);
        });
      }
    });
  }

  public openAddTrigger(): void {
    console.log(this.selection)
    const dialogRef = this.dialog.open(AddFlowDialogComponent, {
      data: {
        flowType: this.flowType,
        triggerType: TriggerType.TRIGGER_TYPE_POST_AUTHENTICATION,
        actions: (this.selection && this.selection.length) ? this.selection : [],
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((req: SetTriggerActionsRequest) => {
      if (req) {
        this.mgmtService.setTriggerActions(req.getActionIdsList(), req.getFlowType(), req.getTriggerType()).then(resp => {
          this.loadFlow();
        }).catch((error: any) => {
          this.toast.showError(error);
        });
      }
    });
  }

}
