import { AfterViewInit, Component, computed, effect, signal, ViewChild } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import { ActionsTwoAddActionTypeComponent } from './actions-two-add-action-type/actions-two-add-action-type.component';
import { MessageInitShape } from '@bufbuild/protobuf';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/resources/action/v3alpha/action_service_pb';
import {
  ActionsTwoAddActionConditionComponent,
  ConditionType,
} from './actions-two-add-action-condition/actions-two-add-action-condition.component';
import {
  ActionsTwoAddActionTargetComponent,
  TargetInit,
} from './actions-two-add-action-target/actions-two-add-action-target.component';
import { CommonModule } from '@angular/common';

enum Page {
  Type,
  Condition,
  Target,
}

type ConditionInit = NonNullable<MessageInitShape<typeof SetExecutionRequestSchema>['condition']>['conditionType'];

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
  public conditionSignal = signal<ConditionInit | undefined>(undefined);
  public targetSignal = signal<TargetInit | undefined>(undefined);

  public request = computed<MessageInitShape<typeof SetExecutionRequestSchema>>(() => {
    return {
      condition: {
        conditionType: this.conditionSignal(),
      },
      execution: {
        targets: [
          {
            type: this.targetSignal(),
          },
        ],
      },
    };
  });

  constructor(public dialogRef: MatDialogRef<ActionTwoAddActionDialogComponent>) {}

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

  public previous() {
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
