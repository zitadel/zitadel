import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { InputModule } from 'src/app/modules/input/input.module';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { EMPTY, Observable, catchError, defer, filter, map, of, shareReplay, startWith, switchMap, tap } from 'rxjs';
import { ExecutionType } from '../actions-two-add-action-type/actions-two-add-action-type.component';
import { MatRadioModule } from '@angular/material/radio';
import {
  EventExecution,
  FunctionExecution,
  RequestExecution,
  ResponseExecution,
} from '@zitadel/proto/zitadel/resources/action/v3alpha/execution_pb';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { F } from '@angular/cdk/keycodes';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';
import { MatCheckboxModule } from '@angular/material/checkbox';

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
  protected readonly conditionForm: ReturnType<typeof this.buildActionConditionFormForRequest> =
    this.buildActionConditionFormForRequest();

  @Output() public continue: EventEmitter<void> = new EventEmitter();
  // @Output() public conditionChanges$: Observable<RequestExecution | ResponseExecution | FunctionExecution | EventExecution>;

  public readonly executionServices$: Observable<string[] | undefined> = of(undefined);
  public readonly executionMethods$: Observable<string[] | undefined> = of(undefined);
  public readonly executionFunctions$: Observable<string[] | undefined> = of(undefined);

  public readonly conditionChanges$: Observable<string[] | undefined> = of(undefined);

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
  }

  public ngOnInit(): void {}

  public buildActionConditionFormForRequest() {
    return this.fb.group({
      all: new FormControl<boolean>(true),
      service: new FormControl<string>(''),
      method: new FormControl<string>(''),
    });
  }

  public buildActionConditionFormForResponse() {
    return this.fb.group({
      all: new FormControl<boolean>(true),
      service: new FormControl<string>(''),
      method: new FormControl<string>(''),
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
