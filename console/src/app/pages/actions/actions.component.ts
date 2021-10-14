import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { Flow, FlowType, TriggerType } from 'src/app/proto/generated/zitadel/action_pb';
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

  public openAddTrigger(): void {
    const dialogRef = this.dialog.open(AddFlowDialogComponent, {
      data: {
        flowType: this.flowType,
        triggerType: TriggerType.TRIGGER_TYPE_POST_AUTHENTICATION,
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
