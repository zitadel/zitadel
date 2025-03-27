import { AfterViewInit, Component, effect, signal, ViewChild } from '@angular/core';
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
import { ActionsTwoAddActionTargetComponent } from './actions-two-add-action-target/actions-two-add-action-target.component';
import { CommonModule } from '@angular/common';
import { BehaviorSubject, combineLatest, forkJoin, map, merge, Observable, of, ReplaySubject } from 'rxjs';
import { FunctionExecution } from '@zitadel/proto/zitadel/resources/action/v3alpha/execution_pb';
import { ActionsTwoTargetsComponent } from '../actions-two-targets/actions-two-targets.component';
import { Condition } from 'src/app/proto/generated/zitadel/action/v2beta/execution_pb';

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
  public conditionSignal = signal<any>(undefined);
  public targetSignal = signal<string>('');

  constructor(public dialogRef: MatDialogRef<ActionTwoAddActionDialogComponent>) {
    effect(() => {
      const type = this.typeSignal();
      const condition = this.conditionSignal();
      const target = this.targetSignal();

      console.log('Request:', type, condition, target);
      const req: MessageInitShape<typeof SetExecutionRequestSchema> = {
        condition: {
          // Map condition here
        },
        execution: {
          // Map execution here
        },
      };
      return req;
      // Perform any additional logic here
      // this.executionType$.next(type); // Example: Update executionType$ with the new value
    });
  }

  public continue() {
    const currentPage = this.page();
    if (currentPage === Page.Type) {
      this.page.set(Page.Condition);
    } else if (currentPage === Page.Condition) {
      this.page.set(Page.Target);
    } else {
      this.dialogRef.close();
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

  public closeWithResult() {
    this.dialogRef.close();
  }

  public onTypeChanged(type: ConditionType): void {
    this.typeSignal.set(type);
  }

  public onConditionChanged(condition: any): void {
    this.conditionSignal.set(condition);
  }

  public onTargetChanged(target: string): void {
    this.targetSignal.set(target);
  }
}
