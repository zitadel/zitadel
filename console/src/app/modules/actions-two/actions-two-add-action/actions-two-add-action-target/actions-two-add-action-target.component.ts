import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Observable, catchError, defer, map, of, shareReplay, ReplaySubject, combineLatestWith } from 'rxjs';
import { MatRadioModule } from '@angular/material/radio';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { InputModule } from 'src/app/modules/input/input.module';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MessageInitShape } from '@bufbuild/protobuf';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/action/v2beta/action_service_pb';
import { Condition, ExecutionTargetTypeSchema } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { MatSelectModule } from '@angular/material/select';
import { atLeastOneFieldValidator } from 'src/app/modules/form-field/validators/validators';
import { ActionConditionPipeModule } from 'src/app/pipes/action-condition-pipe/action-condition-pipe.module';

export type TargetInit = NonNullable<
  NonNullable<MessageInitShape<typeof SetExecutionRequestSchema>['targets']>
>[number]['type'];

@Component({
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  selector: 'cnsl-actions-two-add-action-target',
  templateUrl: './actions-two-add-action-target.component.html',
  styleUrls: ['./actions-two-add-action-target.component.scss'],
  imports: [
    TranslateModule,
    MatRadioModule,
    RouterModule,
    ReactiveFormsModule,
    InputModule,
    MatAutocompleteModule,
    FormsModule,
    ActionConditionPipeModule,
    CommonModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatSelectModule,
  ],
})
export class ActionsTwoAddActionTargetComponent {
  protected readonly targetForm = this.buildActionTargetForm();

  @Output() public readonly back = new EventEmitter<void>();
  @Output() public readonly continue = new EventEmitter<MessageInitShape<typeof ExecutionTargetTypeSchema>[]>();
  @Input() public hideBackButton = false;
  @Input() set selectedCondition(selectedCondition: Condition | undefined) {
    this.selectedCondition$.next(selectedCondition);
  }

  private readonly selectedCondition$ = new ReplaySubject<Condition | undefined>(1);

  protected readonly executionTargets$: Observable<Target[]>;
  protected readonly executionConditions$: Observable<Condition[]>;

  constructor(
    private readonly fb: FormBuilder,
    private readonly actionService: ActionService,
    private readonly toast: ToastService,
  ) {
    this.executionTargets$ = this.listExecutionTargets().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionConditions$ = this.listExecutionConditions().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
  }

  private buildActionTargetForm() {
    return this.fb.group(
      {
        target: new FormControl<Target | null>(null, { validators: [] }),
        executionConditions: new FormControl<Condition[]>([], { validators: [] }),
      },
      {
        validators: atLeastOneFieldValidator(['target', 'executionConditions']),
      },
    );
  }

  private listExecutionTargets() {
    return defer(() => this.actionService.listTargets({})).pipe(
      map(({ result }) => result.filter(this.targetHasDetailsAndConfig)),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private listExecutionConditions(): Observable<Condition[]> {
    const selectedConditionJson$ = this.selectedCondition$.pipe(map((c) => JSON.stringify(c)));

    return defer(() => this.actionService.listExecutions({})).pipe(
      combineLatestWith(selectedConditionJson$),
      map(([executions, selectedConditionJson]) =>
        executions.result.map((e) => e?.condition).filter(this.conditionIsDefinedAndNotCurrentOne(selectedConditionJson)),
      ),

      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private conditionIsDefinedAndNotCurrentOne(selectedConditionJson?: string) {
    return (condition?: Condition): condition is Condition => {
      if (!condition) {
        // condition is undefined so it is not of type Condition
        return false;
      }
      if (!selectedConditionJson) {
        // condition is defined, and we don't have a selectedCondition so we can return all conditions
        return true;
      }
      // we only return conditions that are not the same as the selectedCondition
      return JSON.stringify(condition) !== selectedConditionJson;
    };
  }

  private targetHasDetailsAndConfig(target: Target): target is Target {
    return !!target.id && !!target.id;
  }

  protected submit() {
    const { target, executionConditions } = this.targetForm.getRawValue();

    let valueToEmit: MessageInitShape<typeof ExecutionTargetTypeSchema>[] = target
      ? [
          {
            type: {
              case: 'target',
              value: target.id,
            },
          },
        ]
      : [];

    const includeConditions: MessageInitShape<typeof ExecutionTargetTypeSchema>[] = executionConditions
      ? executionConditions.map((condition) => ({
          type: {
            case: 'include',
            value: condition,
          },
        }))
      : [];

    valueToEmit = [...valueToEmit, ...includeConditions];

    this.continue.emit(valueToEmit);
  }
}
