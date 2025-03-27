import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, DestroyRef, EventEmitter, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { FormBuilder, FormControl, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Observable, catchError, defer, map, of, shareReplay } from 'rxjs';
import { MatRadioModule } from '@angular/material/radio';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { InputModule } from 'src/app/modules/input/input.module';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { GetTarget } from '@zitadel/proto/zitadel/resources/action/v3alpha/target_pb';
import { MessageInitShape } from '@bufbuild/protobuf';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/resources/action/v3alpha/action_service_pb';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { tap } from 'rxjs/operators';

type Target = Required<Pick<GetTarget, 'config' | 'details'>> & GetTarget;
export type TargetInit = NonNullable<
  NonNullable<MessageInitShape<typeof SetExecutionRequestSchema>['execution']>['targets']
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
    CommonModule,
    MatButtonModule,
    MatProgressSpinnerModule,
  ],
})
export class ActionsTwoAddActionTargetComponent {
  protected readonly targetForm = this.buildActionTargetForm();

  @Output() public readonly target = new EventEmitter<TargetInit>();

  protected readonly executionTargets$: Observable<Target[]>;

  constructor(
    private readonly fb: FormBuilder,
    private readonly actionService: ActionService,
    private readonly toast: ToastService,
    destroyRef: DestroyRef,
  ) {
    this.executionTargets$ = this.listExecutionTargets().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.targetForm.valueChanges.pipe(takeUntilDestroyed(destroyRef)).subscribe(() => this.submit());
  }

  public buildActionTargetForm() {
    return this.fb.group({
      target: new FormControl<Target | null>(null, { validators: [Validators.required] }),
    });
  }

  private listExecutionTargets() {
    return defer(() => this.actionService.searchTargets({})).pipe(
      map(({ result }) => result.filter(this.targetHasDetailsAndConfig)),
      tap((target) => console.log('hoi', target)),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
    );
  }

  private targetHasDetailsAndConfig(target: GetTarget): target is Target {
    return !!target.details && !!target.config;
  }

  private submit() {
    const { target } = this.targetForm.getRawValue();
    if (!target) {
      return;
    }
    this.target.emit({
      case: 'target',
      value: target.details.id,
    });
  }

  protected displayTarget(target?: Target) {
    if (!target) {
      return '';
    }
    return target.config.name;
  }
}
