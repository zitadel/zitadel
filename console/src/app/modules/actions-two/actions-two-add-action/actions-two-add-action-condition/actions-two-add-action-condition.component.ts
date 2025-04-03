import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, DestroyRef, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { InputModule } from 'src/app/modules/input/input.module';
import {
  AbstractControl,
  FormBuilder,
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  ValidationErrors,
  ValidatorFn,
} from '@angular/forms';
import {
  Observable,
  catchError,
  defer,
  map,
  of,
  shareReplay,
  ReplaySubject,
  ObservedValueOf,
  switchMap,
  combineLatestWith,
  OperatorFunction,
} from 'rxjs';
import { MatRadioModule } from '@angular/material/radio';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { atLeastOneFieldValidator, requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { Message } from '@bufbuild/protobuf';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { Condition } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { startWith } from 'rxjs/operators';

export type ConditionType = NonNullable<Condition['conditionType']['case']>;
export type ConditionTypeValue<T extends ConditionType> = Omit<
  NonNullable<Extract<Condition['conditionType'], { case: T }>['value']>,
  // we remove the message keys so $typeName is not required
  keyof Message
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
export class ActionsTwoAddActionConditionComponent<T extends ConditionType = ConditionType> {
  @Input({ required: true }) public set conditionType(conditionType: T) {
    this.conditionType$.next(conditionType);
  }
  @Output() public readonly back = new EventEmitter<void>();
  @Output() public readonly continue = new EventEmitter<ConditionTypeValue<T>>();

  private readonly conditionType$ = new ReplaySubject<T>(1);
  protected readonly form$: ReturnType<typeof this.buildForm>;

  protected readonly executionServices$: Observable<string[]>;
  protected readonly executionMethods$: Observable<string[]>;
  protected readonly executionFunctions$: Observable<string[]>;

  constructor(
    private readonly fb: FormBuilder,
    private readonly actionService: ActionService,
    private readonly toast: ToastService,
    private readonly destroyRef: DestroyRef,
  ) {
    this.form$ = this.buildForm().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionServices$ = this.listExecutionServices(this.form$).pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionMethods$ = this.listExecutionMethods(this.form$).pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionFunctions$ = this.listExecutionFunctions(this.form$).pipe(shareReplay({ refCount: true, bufferSize: 1 }));
  }

  public buildForm() {
    return this.conditionType$.pipe(
      switchMap((conditionType) => {
        if (conditionType === 'event') {
          return this.buildEventForm();
        }
        if (conditionType === 'function') {
          return this.buildFunctionForm();
        }
        return this.buildRequestOrResponseForm(conditionType);
      }),
    );
  }

  private buildRequestOrResponseForm<T extends 'request' | 'response'>(requestOrResponse: T) {
    const formFactory = () => ({
      case: requestOrResponse,
      form: this.fb.group(
        {
          all: new FormControl<boolean>(false, { nonNullable: true }),
          service: new FormControl<string>('', { nonNullable: true }),
          method: new FormControl<string>('', { nonNullable: true }),
        },
        {
          validators: atLeastOneFieldValidator(['all', 'service', 'method']),
        },
      ),
    });

    return new Observable<ReturnType<typeof formFactory>>((obs) => {
      const form = formFactory();
      obs.next(form);

      const { all, service, method } = form.form.controls;
      return all.valueChanges
        .pipe(
          map(() => all.value),
          takeUntilDestroyed(this.destroyRef),
        )
        .subscribe((all) => {
          this.toggleFormControl(service, !all);
          this.toggleFormControl(method, !all);
        });
    });
  }

  public buildFunctionForm() {
    return of({
      case: 'function' as const,
      form: this.fb.group({
        name: new FormControl<string>('', { nonNullable: true, validators: [requiredValidator] }),
      }),
    });
  }

  public buildEventForm() {
    const formFactory = () => ({
      case: 'event' as const,
      form: this.fb.group({
        all: new FormControl<boolean>(false, { nonNullable: true }),
        group: new FormControl<string>('', { nonNullable: true }),
        event: new FormControl<string>('', { nonNullable: true }),
      }),
    });

    return new Observable<ReturnType<typeof formFactory>>((obs) => {
      const form = formFactory();
      obs.next(form);

      const { all, group, event } = form.form.controls;
      return all.valueChanges
        .pipe(
          map(() => all.value),
          takeUntilDestroyed(this.destroyRef),
        )
        .subscribe((all) => {
          this.toggleFormControl(group, !all);
          this.toggleFormControl(event, !all);
        });
    });
  }

  private toggleFormControl(control: FormControl, enabled: boolean) {
    if (enabled) {
      control.enable();
    } else {
      control.disable();
    }
  }

  private listExecutionServices(form$: typeof this.form$) {
    return defer(() => this.actionService.listExecutionServices({})).pipe(
      map(({ services }) => services),
      this.formFilter(form$, (form) => {
        if ('service' in form.form.controls) {
          return form.form.controls.service;
        }
        return undefined;
      }),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private listExecutionFunctions(form$: typeof this.form$) {
    return defer(() => this.actionService.listExecutionFunctions({})).pipe(
      map(({ functions }) => functions),
      this.formFilter(form$, (form) => {
        if (form.case !== 'function') {
          return undefined;
        }
        return form.form.controls.name;
      }),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private listExecutionMethods(form$: typeof this.form$) {
    return defer(() => this.actionService.listExecutionMethods({})).pipe(
      map(({ methods }) => methods),
      this.formFilter(form$, (form) => {
        if ('method' in form.form.controls) {
          return form.form.controls.method;
        }
        return undefined;
      }),
      // we also filter by service name
      this.formFilter(form$, (form) => {
        if ('service' in form.form.controls) {
          return form.form.controls.service;
        }
        return undefined;
      }),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private formFilter(
    form$: typeof this.form$,
    getter: (form: ObservedValueOf<typeof this.form$>) => FormControl<string> | undefined,
  ): OperatorFunction<string[], string[]> {
    const filterValue$ = form$.pipe(
      map(getter),
      switchMap((control) => {
        if (!control) {
          return of('');
        }

        return control.valueChanges.pipe(
          startWith(control.value),
          map((value) => value.toLowerCase()),
        );
      }),
    );

    return (obs) =>
      obs.pipe(
        combineLatestWith(filterValue$),
        map(([values, filterValue]) => values.filter((v) => v.toLowerCase().includes(filterValue))),
      );
  }

  protected submit(form: ObservedValueOf<typeof this.form$>) {
    if (form.case === 'request' || form.case === 'response') {
      (this as unknown as ActionsTwoAddActionConditionComponent<'request' | 'response'>).submitRequestOrResponse(form);
    } else if (form.case === 'event') {
      (this as unknown as ActionsTwoAddActionConditionComponent<'event'>).submitEvent(form);
    } else if (form.case === 'function') {
      (this as unknown as ActionsTwoAddActionConditionComponent<'function'>).submitFunction(form);
    }
  }

  private submitRequestOrResponse(
    this: ActionsTwoAddActionConditionComponent<'request' | 'response'>,
    { form }: ObservedValueOf<ReturnType<typeof this.buildRequestOrResponseForm>>,
  ) {
    const { all, service, method } = form.getRawValue();

    if (all) {
      this.continue.emit({
        condition: {
          case: 'all',
          value: true,
        },
      });
    } else if (method) {
      this.continue.emit({
        condition: {
          case: 'method',
          value: method,
        },
      });
    } else if (service) {
      this.continue.emit({
        condition: {
          case: 'service',
          value: service,
        },
      });
    }
  }

  private submitEvent(
    this: ActionsTwoAddActionConditionComponent<'event'>,
    { form }: ObservedValueOf<ReturnType<typeof this.buildEventForm>>,
  ) {
    const { all, event, group } = form.getRawValue();
    if (all) {
      this.continue.emit({
        condition: {
          case: 'all',
          value: true,
        },
      });
    } else if (event) {
      this.continue.emit({
        condition: {
          case: 'event',
          value: event,
        },
      });
    } else if (group) {
      this.continue.emit({
        condition: {
          case: 'group',
          value: group,
        },
      });
    }
  }

  private submitFunction(
    this: ActionsTwoAddActionConditionComponent<'function'>,
    { form }: ObservedValueOf<ReturnType<typeof this.buildFunctionForm>>,
  ) {
    const { name } = form.getRawValue();
    this.continue.emit({
      name,
    });
  }
}
