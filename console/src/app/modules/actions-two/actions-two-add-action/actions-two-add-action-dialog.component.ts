import { Component, signal, ViewChild } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import {
  ActionsTwoAddActionTypeComponent,
  ExecutionType,
} from './actions-two-add-action-type/actions-two-add-action-type.component';
import { MessageInitShape } from '@bufbuild/protobuf';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/resources/action/v3alpha/action_service_pb';
import {
  ActionsTwoAddActionConditionComponent,
  ConditionType,
} from './actions-two-add-action-condition/actions-two-add-action-condition.component';
import { ActionsTwoAddActionTargetComponent } from './actions-two-add-action-target/actions-two-add-action-target.component';
import { CommonModule } from '@angular/common';
import { Observable, of, ReplaySubject } from 'rxjs';
import { FunctionExecution } from '@zitadel/proto/zitadel/resources/action/v3alpha/execution_pb';

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
  protected readonly conditionType: 'function' = 'function';
  public Page = Page;
  @ViewChild('actionTypeComponent', { static: false }) actionTypeComponent!: ActionsTwoAddActionTypeComponent;
  @ViewChild('actionConditionComponent') actionConditionComponent!: ActionsTwoAddActionConditionComponent<ConditionType>;
  @ViewChild('actionTargetComponent') actionTargetComponent!: ActionsTwoAddActionTargetComponent;

  public page = signal<Page | undefined>(Page.Type);
  private request$: Observable<MessageInitShape<typeof SetExecutionRequestSchema>> = of({});

  public executionType$ = new ReplaySubject<ExecutionType>(1);

  constructor(public dialogRef: MatDialogRef<ActionTwoAddActionDialogComponent>) {}

  // ngAfterViewInit(): void {
  // this.actionTypeComponent?.typeChanges$.subscribe((type) => {
  //   console.log('Execution type changed:', type);
  // });
  // this.request$ = forkJoin({
  //   type: this.actionTypeComponent?.typeChanges$,
  //   condition: this.actionConditionComponent.conditionChanges$,
  //   target: this.actionTargetComponent.targetChanges$,
  // }).pipe(
  //   map(({ type, condition, target }) => {
  //     console.log('Request:', type, condition, target);
  //     const req: MessageInitShape<typeof SetExecutionRequestSchema> = {
  //       condition: {
  //         // conditionType: {
  //         // }
  //       },
  //       execution: {},
  //     };
  //     return req;
  //   }),
  // );
  // }

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

  public onTypeChanged(type: ExecutionType): void {
    this.executionType$.next(type);
  }

  public onConditionChanged(condition: any): void {
    console.log('condition changed:', condition);
  }

  public onTargetChanged(target: string): void {
    console.log('target changed:', target);
  }

  public reconstructRequest(): void {}

  protected readonly ExecutionType = ExecutionType;

  test($event: FunctionExecution) {}
}
