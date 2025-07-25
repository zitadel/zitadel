import { CdkDragDrop, moveItemInArray } from '@angular/cdk/drag-drop';
import { Component, DestroyRef } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActionKeysType } from 'src/app/modules/action-keys/action-keys.component';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Action, ActionState, Flow, FlowType, TriggerType } from 'src/app/proto/generated/zitadel/action_pb';
import { SetTriggerActionsRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddFlowDialogComponent } from './add-flow-dialog/add-flow-dialog.component';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'cnsl-actions',
  templateUrl: './actions.component.html',
  styleUrls: ['./actions.component.scss'],
})
export class ActionsComponent {
  protected flow!: Flow.AsObject;

  protected typeControl: UntypedFormControl = new UntypedFormControl();

  protected typesForSelection: FlowType.AsObject[] = [];

  protected selection: Action.AsObject[] = [];
  protected InfoSectionType = InfoSectionType;
  protected ActionKeysType = ActionKeysType;

  protected maxActions: number | null = null;
  protected ActionState = ActionState;
  constructor(
    private mgmtService: ManagementService,
    breadcrumbService: BreadcrumbService,
    private dialog: MatDialog,
    private toast: ToastService,
    destroyRef: DestroyRef,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);

    this.getFlowTypes().then();

    this.typeControl.valueChanges.pipe(takeUntilDestroyed(destroyRef)).subscribe((value) => {
      this.loadFlow((value as FlowType.AsObject).id);
    });
  }

  private async getFlowTypes(): Promise<void> {
    try {
      let resp = await this.mgmtService.listFlowTypes();
      this.typesForSelection = resp.resultList;
      if (!this.flow && resp.resultList[0]) {
        const type = resp.resultList[0];
        this.typeControl.setValue(type);
      }
    } catch (error) {
      this.toast.showError(error);
    }
  }

  private loadFlow(id: string) {
    this.mgmtService.getFlow(id).then((flowResponse) => {
      if (flowResponse.flow) {
        this.flow = flowResponse.flow;
      }
    });
  }

  public clearFlow(id: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.CLEAR',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'FLOWS.DIALOG.CLEAR.TITLE',
        descriptionKey: 'FLOWS.DIALOG.CLEAR.DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtService
          .clearFlow(id)
          .then(() => {
            this.toast.showInfo('FLOWS.FLOWCLEARED', true);
            this.loadFlow(id);
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    });
  }

  protected openAddTrigger(flow: FlowType.AsObject, trigger?: TriggerType.AsObject): void {
    const dialogRef = this.dialog.open(AddFlowDialogComponent, {
      data: {
        flowType: flow,
        actions: this.selection && this.selection.length ? this.selection : [],
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((req: SetTriggerActionsRequest) => {
      if (req) {
        this.mgmtService
          .setTriggerActions(req.getActionIdsList(), req.getFlowType(), req.getTriggerType())
          .then(() => {
            this.toast.showInfo('FLOWS.FLOWCHANGED', true);
            this.loadFlow(flow.id);
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    });
  }

  drop(triggerActionsListIndex: number, array: any[], event: CdkDragDrop<Action.AsObject[]>) {
    moveItemInArray(array, event.previousIndex, event.currentIndex);
    this.saveFlow(triggerActionsListIndex);
  }

  saveFlow(index: number) {
    if (
      this.flow.type &&
      this.flow.triggerActionsList &&
      this.flow.triggerActionsList[index] &&
      this.flow.triggerActionsList[index]?.triggerType
    ) {
      this.mgmtService
        .setTriggerActions(
          this.flow.triggerActionsList[index].actionsList.map((action) => action.id),
          this.flow.type.id,
          this.flow.triggerActionsList[index].triggerType?.id ?? '',
        )
        .then(() => {
          this.toast.showInfo('FLOWS.TOAST.ACTIONSSET', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  protected removeTriggerActionsList(index: number) {
    if (this.flow.type && this.flow.triggerActionsList && this.flow.triggerActionsList[index]) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.CLEAR',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'FLOWS.DIALOG.REMOVEACTIONSLIST.TITLE',
          descriptionKey: 'FLOWS.DIALOG.REMOVEACTIONSLIST.DESCRIPTION',
        },
        width: '400px',
      });

      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.mgmtService
            .setTriggerActions([], this.flow?.type?.id ?? '', this.flow.triggerActionsList[index].triggerType?.id ?? '')
            .then(() => {
              this.toast.showInfo('FLOWS.TOAST.ACTIONSSET', true);
              this.loadFlow(this.flow?.type?.id ?? '');
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    }
  }
}
