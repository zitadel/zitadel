import { AfterViewInit, Component, signal, ViewChild } from '@angular/core';
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
import { combineLatest, forkJoin, map, merge, Observable, of, ReplaySubject } from 'rxjs';
import { FunctionExecution } from '@zitadel/proto/zitadel/resources/action/v3alpha/execution_pb';
import { ActionsTwoTargetsComponent } from '../actions-two-targets/actions-two-targets.component';

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
export class ActionTwoAddActionDialogComponent implements AfterViewInit {
  protected readonly conditionType: ConditionType = 'function';
  public Page = Page;
  @ViewChild(ActionsTwoAddActionTypeComponent, { static: false }) actionTypeComponent!: ActionsTwoAddActionTypeComponent;
  @ViewChild(ActionsTwoAddActionConditionComponent, { static: false })
  actionConditionComponent!: ActionsTwoAddActionConditionComponent<ConditionType>;
  @ViewChild(ActionsTwoAddActionTargetComponent, { static: false })
  actionTargetComponent!: ActionsTwoAddActionTargetComponent;

  public page = signal<Page | undefined>(Page.Type);
  private request$: Observable<MessageInitShape<typeof SetExecutionRequestSchema>> = of({});

  public executionType$ = new ReplaySubject<ExecutionType>(1);

  private typeState$ = new ReplaySubject<ExecutionType | null>(1);
  private conditionState$ = new ReplaySubject<any | null>(1);
  private targetState$ = new ReplaySubject<string | null>(1);

  constructor(public dialogRef: MatDialogRef<ActionTwoAddActionDialogComponent>) {
    this.typeState$.subscribe((value) => console.log('Type$:', value));
    this.conditionState$.subscribe((value) => console.log('Condition$:', value));
    this.targetState$.subscribe((value) => console.log('Target$:', value));

    // Combine the ReplaySubjects into a single Observable
    this.request$ = combineLatest({
      type: this.typeState$,
      condition: this.conditionState$,
      target: this.targetState$,
    }).pipe(
      map(({ type, condition, target }) => {
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
      }),
    );
  }

  ngAfterViewInit(): void {
    // Pipe the Observables to the ReplaySubjects cause the ViewChilds are not available all the time and merge() does not work
    this.actionTypeComponent?.typeChanges$.subscribe(this.typeState$);
    this.actionConditionComponent?.conditionTypeValue.subscribe(this.conditionState$);
    this.actionTargetComponent?.targetChanges$.subscribe(this.targetState$);
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
