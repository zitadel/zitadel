import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { InputModule } from 'src/app/modules/input/input.module';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { EMPTY, Observable, catchError, defer, filter, map, of, shareReplay, startWith, switchMap, tap } from 'rxjs';
import { ExecutionType } from '../actions-two-add-action-type/actions-two-add-action-type.component';
import { MatRadioModule } from '@angular/material/radio';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { Condition } from '@zitadel/proto/zitadel/resources/action/v3alpha/execution_pb';

@Component({
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  selector: 'cnsl-actions-two-add-action-condition',
  templateUrl: './actions-two-add-action-condition.component.html',
  styleUrls: ['./actions-two-add-action-condition.component.scss'],
  imports: [
    TranslateModule,
    MatRadioModule,
    RouterModule,
    ReactiveFormsModule,
    InputModule,
    MatAutocompleteModule,
    MatCheckboxModule,
    FormsModule,
    CommonModule,
    MatButtonModule,
    MatProgressSpinnerModule,
  ],
})
export class ActionsTwoAddActionConditionComponent implements OnInit {
  public ExecutionType = ExecutionType;
  protected conditionForm:
    | ReturnType<typeof this.buildActionConditionFormForRequestOrResponse>
    | ReturnType<typeof this.buildActionConditionFormForFunctions>
    | ReturnType<typeof this.buildActionConditionFormForEvents> = this.buildActionConditionFormForRequestOrResponse();

  @Output() public continue: EventEmitter<void> = new EventEmitter();
  // @Output() public conditionChanges$: Observable<RequestExecution | ResponseExecution | FunctionExecution | EventExecution>;

  public readonly executionServices$: Observable<string[] | undefined> = of(undefined);
  public readonly executionMethods$: Observable<string[] | undefined> = of(undefined);
  public readonly executionFunctions$: Observable<string[] | undefined> = of(undefined);

  @Output() public conditionChanges$!: Observable<Condition>;
  @Input() public executionType$!: Observable<ExecutionType>;

  constructor(
    private readonly fb: FormBuilder,
    private actionService: ActionService,
    private toast: ToastService,
  ) {
    // Initialize the Observable to emit form value changes
    // this.conditionChanges$ = this.conditionForm!.valueChanges.pipe(
    //   startWith(this.conditionForm!.value), // Emit the initial value
    //   tap((value) => console.log('ExecutionType changed:', value)), // Debugging/logging
    // );
    this.executionServices$ = this.listExecutionServices().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionMethods$ = this.listExecutionMethods().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionFunctions$ = this.listExecutionFunctions().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    // const condition: MessageInitShape<typeof ConditionSchema> = {
    //   conditionType: {
    //     case: 'request',
    //     value: {
    //       condition: this.conditionForm.get('all').value
    //         ? {
    //             case: 'all',
    //           }
    //         : this.conditionForm.get('service')?.value
    //           ? {
    //               case: 'service',
    //               value: {
    //                 service: this.conditionForm.get('service')?.value,
    //                 method: this.conditionForm.get('method')?.value,
    //               },
    //             }
    //           : this.conditionForm.get('method')?.value
    //             ? {
    //                 case: 'method',
    //                 value: {
    //                   method: this.conditionForm.get('method')?.value,
    //                 },
    //               }
    //             : undefined,
    //     },
    //   },
    // };

    const r = this.conditionForm.valueChanges.pipe();
  }

  public ngOnInit(): void {
    // Subscribe to executionType$ to get the latest value
    this.executionType$
      .pipe(
        tap((type) => console.log('ExecutionType received in condition component:', type)),
        map((type) => {
          // Dynamically update the form based on the execution type
          switch (type) {
            case ExecutionType.EVENTS:
              return this.buildActionConditionFormForEvents();
            case ExecutionType.FUNCTIONS:
              return this.buildActionConditionFormForFunctions();
            default:
              return this.buildActionConditionFormForRequestOrResponse();
          }
        }),
      )
      .subscribe((form) => {
        this.conditionForm = form; // Update the form dynamically
      });

    // // @ts-ignore
    // this.conditionChanges$ = this.executionType$.pipe(
    //   switchMap((executionType) =>
    //     this.conditionForm.valueChanges.pipe(
    //       //@ts-ignore
    //       startWith(this.conditionForm.value),
    //       map((formValues) => this.mapToCondition(executionType, formValues)),
    //     ),
    //   ),
    //   tap((condition) => console.log('Mapped Condition:', condition)), // Debugging
    // );
  }

  private mapToCondition(executionType: ExecutionType, formValues: any): Condition {
    type ExecutionTypeValue = (typeof ExecutionType)[keyof typeof ExecutionType];

    switch (executionType) {
      case ExecutionType.FUNCTIONS:
        return {
          conditionType: {
            case: 'function',
            value: {
              name: formValues.name,
            },
          },
        } as Condition;

      case ExecutionType.EVENTS:
        const conditionValue: any = formValues.all
          ? { condition: { case: 'all', value: true } }
          : formValues.service
            ? {
                condition: {
                  case: 'service',
                  value: formValues.service as string,
                },
              }
            : formValues.method
              ? {
                  condition: {
                    case: 'method',
                    value: formValues.method as string,
                  },
                }
              : undefined;

        return {
          conditionType: {
            case: executionType as ExecutionTypeValue,
            value: conditionValue,
          },
        } as Condition;

      default:
        const defaultConditionValue = formValues.all
          ? { condition: { case: 'all', value: true } }
          : formValues.service
            ? {
                condition: {
                  case: 'service',
                  value: formValues.service as string,
                },
              }
            : formValues.method
              ? {
                  condition: {
                    case: 'method',
                    value: formValues.method as string,
                  },
                }
              : undefined;

        return {
          conditionType: {
            case: executionType as ExecutionTypeValue,
            value: defaultConditionValue,
          },
        } as Condition;
    }
  }

  public buildActionConditionFormForRequestOrResponse() {
    return this.fb.group({
      all: new FormControl<boolean>(true),
      service: new FormControl<string>(''),
      method: new FormControl<string>(''),
    });
  }

  public buildActionConditionFormForFunctions() {
    return this.fb.group({
      name: new FormControl<string>('', [requiredValidator]),
    });
  }

  public buildActionConditionFormForEvents() {
    return this.fb.group({
      event: new FormControl<string>('', [requiredValidator]),
      group: new FormControl<string>('', [requiredValidator]),
      all: new FormControl<boolean>(true),
    });
  }

  private listExecutionServices() {
    return defer(() => this.actionService.listExecutionServices()).pipe(
      map(({ services }) => services),
      catchError((error) => {
        this.toast.showError(error);
        return EMPTY;
      }),
    );
  }

  private listExecutionFunctions() {
    return defer(() => this.actionService.listExecutionFunctions()).pipe(
      map(({ functions }) => functions),
      catchError((error) => {
        this.toast.showError(error);
        return EMPTY;
      }),
    );
  }

  private listExecutionMethods() {
    return defer(() => this.actionService.listExecutionMethods()).pipe(
      map(({ methods }) => methods),
      catchError((error) => {
        this.toast.showError(error);
        return EMPTY;
      }),
    );
  }
}
