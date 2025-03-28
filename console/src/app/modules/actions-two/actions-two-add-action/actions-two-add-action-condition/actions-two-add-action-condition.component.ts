import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, DestroyRef, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { InputModule } from 'src/app/modules/input/input.module';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
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
  tap,
  Subject,
} from 'rxjs';
import { MatRadioModule } from '@angular/material/radio';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { Message } from '@bufbuild/protobuf';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { Condition } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';

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
export class ActionsTwoAddActionConditionComponent<T extends ConditionType = ConditionType> implements OnInit {
  @Input({ required: true }) public set conditionType(conditionType: T) {
    this.conditionType$.next(conditionType);
  }
  @Output() public readonly conditionTypeValue = new EventEmitter<ConditionTypeValue<T>>();
  @Input() public continue: Subject<void> = new Subject<void>();
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
    this.executionServices$ = this.listExecutionServices().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionMethods$ = this.listExecutionMethods().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executionFunctions$ = this.listExecutionFunctions().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.form$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((form) => this.submit(form));
  }

  ngOnInit(): void {
    this.continue
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        switchMap(() => this.form$),
      )
      .subscribe((form) => {
        this.submit(form);
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
      form: this.fb.group({
        all: new FormControl<boolean>(false, { nonNullable: true }),
        service: new FormControl<string>('', { nonNullable: true }),
        method: new FormControl<string>('', { nonNullable: true }),
      }),
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

  private listExecutionServices() {
    return defer(() => this.actionService.listExecutionServices({})).pipe(
      map(({ services }) => services),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private listExecutionFunctions() {
    return defer(() => this.actionService.listExecutionFunctions({})).pipe(
      map(({ functions }) => functions),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private listExecutionMethods() {
    return defer(() => this.actionService.listExecutionMethods({})).pipe(
      map(({ methods }) => methods),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
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
    console.log(all, service, method);

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
    { form }: ObservedValueOf<ReturnType<typeof this.buildEventForm>>,
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
    { form }: ObservedValueOf<ReturnType<typeof this.buildFunctionForm>>,
  ) {
    const { name } = form.getRawValue();
    this.conditionTypeValue.emit({
      name,
    });
  }
}
