import { Component, computed, effect, Inject, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import { ActionsTwoAddActionTypeComponent } from './actions-two-add-action-type/actions-two-add-action-type.component';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  ActionsTwoAddActionConditionComponent,
  ConditionType,
  ConditionTypeValue,
} from './actions-two-add-action-condition/actions-two-add-action-condition.component';
import { ActionsTwoAddActionTargetComponent } from './actions-two-add-action-target/actions-two-add-action-target.component';
import { CommonModule } from '@angular/common';
import { Execution, ExecutionTargetTypeSchema } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { Subject } from 'rxjs';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/action/v2beta/action_service_pb';

enum Page {
  Type,
  Condition,
  Target,
}

@Component({
  selector: 'cnsl-actions-two-add-action-dialog',
  templateUrl: './actions-two-add-action-dialog.component.html',
  styleUrls: ['./actions-two-add-action-dialog.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatDialogModule,
    TranslateModule,
    ActionsTwoAddActionTypeComponent,
    ActionsTwoAddActionConditionComponent,
    ActionsTwoAddActionTargetComponent,
  ],
})
export class ActionTwoAddActionDialogComponent {
  public Page = Page;
  public page = signal<Page | undefined>(Page.Type);

  public typeSignal = signal<ConditionType>('request');
  public conditionSignal = signal<ConditionTypeValue<ConditionType> | undefined>(undefined); // TODO: fix this type
  public targetSignal = signal<Array<MessageInitShape<typeof ExecutionTargetTypeSchema>> | undefined>(undefined);

  public continueSubject = new Subject<void>();

  public request = computed<MessageInitShape<typeof SetExecutionRequestSchema>>(() => {
    const req = {
      condition: {
        conditionType: {
          case: this.typeSignal(),
          value: this.conditionSignal() as any, // TODO: fix this type
        },
      },
      targets: this.targetSignal(),
    };

    console.log(req);
    return req;
  });

  constructor(
    public dialogRef: MatDialogRef<ActionTwoAddActionDialogComponent, MessageInitShape<typeof SetExecutionRequestSchema>>,
    @Inject(MAT_DIALOG_DATA) protected readonly data: { execution?: Execution },
  ) {
    if (data?.execution) {
      this.typeSignal.set(data.execution.condition?.conditionType.case ?? 'request');
      this.conditionSignal.set((data.execution.condition?.conditionType as any)?.value ?? undefined);
      this.targetSignal.set(data.execution.targets ?? []);

      this.page.set(Page.Target); // Set the initial page based on the provided execution data
    }

    effect(() => {
      const currentPage = this.page();
      if (currentPage === Page.Target) {
        this.continueSubject.next(); // Trigger the Subject to request condition form when the page changes to "Target"
      }
    });
  }

  public continue() {
    const currentPage = this.page();
    if (currentPage === Page.Type) {
      this.page.set(Page.Condition);
    } else if (currentPage === Page.Condition) {
      this.page.set(Page.Target);
    } else {
      this.dialogRef.close(this.request());
    }
  }

  public back() {
    const currentPage = this.page();
    if (currentPage === Page.Target) {
      this.page.set(Page.Condition);
    } else if (currentPage === Page.Condition) {
      this.page.set(Page.Type);
    } else {
      this.dialogRef.close();
    }
  }
}
