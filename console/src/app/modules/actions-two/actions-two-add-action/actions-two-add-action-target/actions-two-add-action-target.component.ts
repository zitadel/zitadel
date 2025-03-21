import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { EMPTY, Observable, catchError, defer, map, of, shareReplay } from 'rxjs';
import { ExecutionType } from '../actions-two-add-action-type/actions-two-add-action-type.component';
import { MatRadioModule } from '@angular/material/radio';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { InputModule } from 'src/app/modules/input/input.module';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

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
    CommonModule,
    MatButtonModule,
    MatProgressSpinnerModule,
  ],
})
export class ActionsTwoAddActionTargetComponent implements OnInit {
  public ExecutionType = ExecutionType;
  protected readonly targetForm: ReturnType<typeof this.buildActionTargetForm> = this.buildActionTargetForm();

  @Output() public continue: EventEmitter<void> = new EventEmitter();
  // @Output() public conditionChanges$: Observable<RequestExecution | ResponseExecution | FunctionExecution | EventExecution>;

  public readonly executionTargets$: Observable<string[] | undefined> = of(undefined);

  constructor(
    private readonly fb: FormBuilder,
    private actionService: ActionService,
    private toast: ToastService,
  ) {
    this.executionTargets$ = this.listExecutionTargets().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
  }

  public ngOnInit(): void {}

  public buildActionTargetForm() {
    return this.fb.group({
      all: new FormControl<boolean>(true),
      service: new FormControl<string>(''),
      method: new FormControl<string>(''),
    });
  }

  private listExecutionTargets() {
    return defer(() => this.actionService.listExecutionFunctions()).pipe(
      map(({ functions }) => functions),
      catchError((error) => {
        this.toast.showError(error);
        return EMPTY;
      }),
    );
  }
}
