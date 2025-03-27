import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { InputModule } from 'src/app/modules/input/input.module';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Observable, catchError, defer, map, of, shareReplay, ReplaySubject, ObservedValueOf, switchMap } from 'rxjs';
import { MatRadioModule } from '@angular/material/radio';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { Condition } from '@zitadel/proto/zitadel/resources/action/v3alpha/execution_pb';
import { Message } from '@bufbuild/protobuf';

export type ConditionType = NonNullable<Condition['conditionType']['case']>;
export type ConditionTypeValue<T extends ConditionType> = Omit<
  NonNullable<Extract<Condition['conditionType'], { case: T }>['value']>,
  // we remove the message keys so $typeName is not required
  keyof Message
>;

type Form = Observable<
  | {
      case: 'request' | 'response';
      form: ReturnType<
        ActionsTwoAddActionConditionComponent<'request' | 'response'>['buildActionConditionFormForRequestOrResponse']
      >;
    }
  | {
      case: 'event';
      form: ReturnType<ActionsTwoAddActionConditionComponent<'event'>['buildActionConditionFormForEvents']>;
    }
  | {
      case: 'function';
      form: ReturnType<ActionsTwoAddActionConditionComponent<'event'>['buildActionConditionFormForFunctions']>;
    }
>;

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
export class ActionsTwoAddActionConditionComponent<T extends ConditionType> {
  @Input({ required: true }) public set conditionType(conditionType: T) {
    this.conditionType$.next(conditionType);
  }
  @Output() public conditionTypeValue = new EventEmitter<ConditionTypeValue<T>>();

  private readonly conditionType$ = new ReplaySubject<T>(1);
  protected readonly form$: Form;

  protected readonly executionServices$: Observable<string[]>;
  protected readonly executionMethods$: Observable<string[]>;
  protected readonly executionFunctions$: Observable<string[]>;

  constructor(
    private readonly fb: FormBuilder,
    private actionService: ActionService,
    private toast: ToastService,
  ) {
    this.form$ = this.buildForm();
    this.executionServices$ = this.listExecutionServices().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionMethods$ = this.listExecutionMethods().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionFunctions$ = this.listExecutionFunctions().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.form$
      .pipe(
        switchMap((form) => {
          console.log(form);
          return form.form.valueChanges;
        }),
      )
      .subscribe(console.log);
  }

  public buildForm(): Form {
    return this.conditionType$.pipe(
      map((conditionType) => {
        if (conditionType === 'request') {
          return {
            case: 'request',
            form: this.buildActionConditionFormForRequestOrResponse(),
          };
        }
        if (conditionType === 'response') {
          return {
            case: 'response',
            form: this.buildActionConditionFormForRequestOrResponse(),
          };
        }
        if (conditionType === 'event') {
          return {
            case: 'event',
            form: this.buildActionConditionFormForEvents(),
          };
        }
        if (conditionType === 'function') {
          return {
            case: 'function',
            form: this.buildActionConditionFormForFunctions(),
          };
        }

        throw new Error('Invalid conditionType');
      }),
    );
  }

  private buildActionConditionFormForRequestOrResponse() {
    return this.fb.group({
      all: new FormControl<boolean>(true, { nonNullable: true }),
      service: new FormControl<string>('', { nonNullable: true }),
      method: new FormControl<string>('', { nonNullable: true }),
    });
  }

  public buildActionConditionFormForFunctions() {
    return this.fb.group({
      name: new FormControl<string>('', { nonNullable: true, validators: [requiredValidator] }),
    });
  }

  public buildActionConditionFormForEvents() {
    return this.fb.group({
      all: new FormControl<boolean>(true, { nonNullable: true }),
      group: new FormControl<string>('', { nonNullable: true }),
      event: new FormControl<string>('', { nonNullable: true }),
    });
  }

  private listExecutionServices() {
    return defer(() => this.actionService.listExecutionServices()).pipe(
      map(({ services }) => services),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private listExecutionFunctions() {
    return defer(() => this.actionService.listExecutionFunctions()).pipe(
      map(({ functions }) => functions),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private listExecutionMethods() {
    return defer(() => this.actionService.listExecutionMethods()).pipe(
      map(({ methods }) => methods),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  protected submit(form: ObservedValueOf<Form>) {
    if (form.case === 'request' || form.case === 'response') {
      (this as unknown as ActionsTwoAddActionConditionComponent<'request' | 'response'>).submitRequestOrResponse(form.form);
    } else if (form.case === 'event') {
      (this as unknown as ActionsTwoAddActionConditionComponent<'event'>).submitEvent(form.form);
    } else if (form.case === 'function') {
      (this as unknown as ActionsTwoAddActionConditionComponent<'function'>).submitFunction(form.form);
    }
  }

  private submitRequestOrResponse(
    this: ActionsTwoAddActionConditionComponent<'request' | 'response'>,
    form: ReturnType<typeof this.buildActionConditionFormForRequestOrResponse>,
  ) {
    const { all, service, method } = form.getRawValue();
    if (all) {
      this.conditionTypeValue.emit({
        condition: {
          case: 'all',
          value: true,
        },
      });
    } else if (method) {
      this.conditionTypeValue.emit({
        condition: {
          case: 'method',
          value: method,
        },
      });
    } else if (service) {
      this.conditionTypeValue.emit({
        condition: {
          case: 'service',
          value: service,
        },
      });
    }
  }

  private submitEvent(
    this: ActionsTwoAddActionConditionComponent<'event'>,
    form: ReturnType<typeof this.buildActionConditionFormForEvents>,
  ) {
    const { all, event, group } = form.getRawValue();
    if (all) {
      this.conditionTypeValue.emit({
        condition: {
          case: 'all',
          value: true,
        },
      });
    } else if (event) {
      this.conditionTypeValue.emit({
        condition: {
          case: 'event',
          value: event,
        },
      });
    } else if (group) {
      this.conditionTypeValue.emit({
        condition: {
          case: 'group',
          value: group,
        },
      });
    }
  }

  private submitFunction(
    this: ActionsTwoAddActionConditionComponent<'function'>,
    form: ReturnType<typeof this.buildActionConditionFormForFunctions>,
  ) {
    const { name } = form.getRawValue();
    this.conditionTypeValue.emit({
      name,
    });
  }
}
