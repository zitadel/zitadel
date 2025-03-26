import { AfterViewInit, Component, Inject, OnChanges, signal, SimpleChanges, ViewChild } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogClose,
  MatDialogContent,
  MatDialogModule,
  MatDialogRef,
  MatDialogTitle,
} from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import {
  ActionsTwoAddActionTypeComponent,
  ExecutionType,
} from './actions-two-add-action-type/actions-two-add-action-type.component';
import { MessageInitShape } from '@bufbuild/protobuf';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/resources/action/v3alpha/action_service_pb';
import { ActionsTwoAddActionConditionComponent } from './actions-two-add-action-condition/actions-two-add-action-condition.component';
import { ActionsTwoAddActionTargetComponent } from './actions-two-add-action-target/actions-two-add-action-target.component';
import { CommonModule } from '@angular/common';
import { forkJoin, map, merge, Observable, of, ReplaySubject } from 'rxjs';

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
  public Page = Page;
  @ViewChild('actionTypeComponent', { static: false }) actionTypeComponent!: ActionsTwoAddActionTypeComponent;
  @ViewChild('actionConditionComponent') actionConditionComponent!: ActionsTwoAddActionConditionComponent;
  @ViewChild('actionTargetComponent') actionTargetComponent!: ActionsTwoAddActionTargetComponent;

  public page = signal<Page | undefined>(Page.Type);
  private request$: Observable<MessageInitShape<typeof SetExecutionRequestSchema>> = of({});

  public executionType$ = new ReplaySubject<ExecutionType>(1);

  constructor(
    public dialogRef: MatDialogRef<ActionTwoAddActionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  ngAfterViewInit(): void {
    // this.actionTypeComponent?.typeChanges$.subscribe((type) => {
    //   console.log('Execution type changed:', type);
    // });

    this.request$ = forkJoin({
      type: this.actionTypeComponent?.typeChanges$,
      condition: this.actionConditionComponent.conditionChanges$,
      target: this.actionTargetComponent.targetChanges$,
    }).pipe(
      map(({ type, condition, target }) => {
        console.log('Request:', type, condition, target);
        const req: MessageInitShape<typeof SetExecutionRequestSchema> = {
          condition: {
            // conditionType: {
            // }
          },
          execution: {},
        };
        return req;
      }),
    );
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
}
