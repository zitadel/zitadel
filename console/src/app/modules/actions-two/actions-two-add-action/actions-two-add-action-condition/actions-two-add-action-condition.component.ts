import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, effect, EventEmitter, Input, Output, Signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { InputModule } from 'src/app/modules/input/input.module';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Observable, of, shareReplay, ReplaySubject, ObservedValueOf, switchMap, Subject } from 'rxjs';
import { MatRadioModule } from '@angular/material/radio';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { atLeastOneFieldValidator, requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { Message } from '@bufbuild/protobuf';
import { Condition } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { distinctUntilChanged, startWith, takeUntil } from 'rxjs/operators';
import { CreateQueryResult } from '@tanstack/angular-query-experimental';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActionsTwoAddActionAutocompleteInputComponent } from '../actions-two-add-action-autocomplete-input/actions-two-add-action-autocomplete-input.component';

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
    ActionsTwoAddActionAutocompleteInputComponent,
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

  protected readonly listExecutionServicesQuery = this.actionService.listExecutionServicesQuery();
  private readonly listExecutionMethodsQuery = this.actionService.listExecutionMethodsQuery();
  protected readonly listExecutionFunctionsQuery = this.actionService.listExecutionFunctionsQuery();
  protected readonly filteredExecutionMethods: Signal<string[] | undefined>;

  constructor(
    private readonly fb: FormBuilder,
    private readonly actionService: ActionService,
    private readonly toast: ToastService,
  ) {
    this.form$ = this.buildForm().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.handleError(this.listExecutionServicesQuery);
    this.handleError(this.listExecutionMethodsQuery);
    this.handleError(this.listExecutionFunctionsQuery);

    this.filteredExecutionMethods = this.filterExecutionMethods(this.form$);
  }

  public handleError(query: CreateQueryResult<any>) {
    return effect(() => {
      const error = query.error();
      if (error) {
        this.toast.showError(error);
      }
    });
  }

  private filterExecutionMethods(form$: typeof this.form$) {
    const service$ = form$.pipe(
      switchMap((form) => {
        if (!('service' in form.form.controls)) {
          return of<string>('');
        }
        const { service } = form.form.controls;
        return service.valueChanges.pipe(startWith(service.value));
      }),
    );

    const serviceSignal = toSignal(service$, { initialValue: '' });

    const query = this.actionService.listExecutionMethodsQuery();

    return computed(() => {
      const methods = query.data();
      const service = serviceSignal();

      if (!methods) {
        return undefined;
      }

      return methods.filter((method) => method.includes(service)).map((method) => method.replace(`/${service}/`, ''));
    });
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

      const destroy$ = new Subject<void>();
      form.form.valueChanges
        .pipe(distinctUntilChanged(undefined!, JSON.stringify), startWith(undefined), takeUntil(destroy$))
        .subscribe(() => {
          this.toggleFormControl(service, !all.value);
          this.toggleFormControl(method, !!service.value && !all.value);
        });

      service.valueChanges.pipe(distinctUntilChanged(), takeUntil(destroy$)).subscribe(() => {
        method.setValue('');
      });

      return () => destroy$.next();
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
      return all.valueChanges.subscribe(() => {
        this.toggleFormControl(group, !all.value);
        this.toggleFormControl(event, !all.value);
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
          value: `/${service}/${method}`,
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
