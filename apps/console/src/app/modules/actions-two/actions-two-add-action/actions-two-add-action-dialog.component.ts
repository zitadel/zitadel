import { Component, computed, effect, Inject, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import { ActionsTwoAddActionTypeComponent } from './actions-two-add-action-type/actions-two-add-action-type.component';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  ActionsTwoAddActionConditionComponent,
  ConditionType,
} from './actions-two-add-action-condition/actions-two-add-action-condition.component';
import { ActionsTwoAddActionTargetComponent } from './actions-two-add-action-target/actions-two-add-action-target.component';
import { CommonModule } from '@angular/common';
import { Condition, Execution } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { Subject } from 'rxjs';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/action/v2beta/action_service_pb';

enum Page {
  Type,
  Condition,
  Target,
}

export type CorrectlyTypedCondition = Condition & { conditionType: Extract<Condition['conditionType'], { case: string }> };

export type CorrectlyTypedExecution = Omit<Execution, 'condition'> & {
  condition: CorrectlyTypedCondition;
};

export const correctlyTypeExecution = (execution: Execution): CorrectlyTypedExecution => {
  if (!execution.condition?.conditionType?.case) {
    throw new Error('Condition is required');
  }
  const conditionType = execution.condition.conditionType;

  const condition = {
    ...execution.condition,
    conditionType,
  };

  return {
    ...execution,
    condition,
  };
};

export type ActionTwoAddActionDialogData = {
  execution?: CorrectlyTypedExecution;
};

export type ActionTwoAddActionDialogResult = MessageInitShape<typeof SetExecutionRequestSchema>;

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
  protected readonly Page = Page;
  protected readonly page = signal<Page>(Page.Type);

  protected readonly typeSignal = signal<ConditionType>('request');
  protected readonly conditionSignal = signal<MessageInitShape<typeof SetExecutionRequestSchema>['condition']>(undefined);
  protected readonly targetsSignal = signal<string[]>([]);

  protected readonly continueSubject = new Subject<void>();

  protected readonly request = computed<MessageInitShape<typeof SetExecutionRequestSchema>>(() => {
    return {
      condition: this.conditionSignal(),
      targets: this.targetsSignal(),
    };
  });

  protected readonly preselectedTargetIds: string[] = [];

  constructor(
    protected readonly dialogRef: MatDialogRef<ActionTwoAddActionDialogComponent, ActionTwoAddActionDialogResult>,
    @Inject(MAT_DIALOG_DATA) protected readonly data: ActionTwoAddActionDialogData,
  ) {
    effect(() => {
      const currentPage = this.page();
      if (currentPage === Page.Target) {
        this.continueSubject.next(); // Trigger the Subject to request condition form when the page changes to "Target"
      }
    });

    if (!data?.execution) {
      return;
    }

    this.targetsSignal.set(data.execution.targets);
    this.typeSignal.set(data.execution.condition.conditionType.case);
    this.conditionSignal.set(data.execution.condition);
    this.preselectedTargetIds = data.execution.targets;

    this.page.set(Page.Target); // Set the initial page based on the provided execution data
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
