import { ChangeDetectionStrategy, Component, DestroyRef, EventEmitter, Input, Output } from '@angular/core';
import { FeatureToggleComponent } from '../feature-toggle.component';
import { ToggleStates } from 'src/app/components/features/features.component';
import { distinctUntilKeyChanged, ReplaySubject } from 'rxjs';
import { FormControl, ReactiveFormsModule, Validators } from '@angular/forms';
import { AsyncPipe, NgIf } from '@angular/common';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';
import { MatTooltipModule } from '@angular/material/tooltip';

@Component({
  standalone: true,
  selector: 'cnsl-login-v2-feature-toggle',
  templateUrl: './login-v2-feature-toggle.component.html',
  imports: [
    FeatureToggleComponent,
    AsyncPipe,
    NgIf,
    ReactiveFormsModule,
    InputModule,
    HasRolePipeModule,
    MatButtonModule,
    TranslateModule,
    MatTooltipModule,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginV2FeatureToggleComponent {
  @Input({ required: true })
  set toggleState(toggleState: ToggleStates['loginV2']) {
    this.toggleState$.next(toggleState);
  }
  @Output()
  public toggleChanged = new EventEmitter<ToggleStates['loginV2']>();

  protected readonly toggleState$ = new ReplaySubject<ToggleStates['loginV2']>(1);
  protected readonly baseUri = new FormControl('', { nonNullable: true, validators: [Validators.required] });

  constructor(destroyRef: DestroyRef) {
    this.toggleState$.pipe(distinctUntilKeyChanged('baseUri'), takeUntilDestroyed(destroyRef)).subscribe(({ baseUri }) => {
      this.baseUri.setValue(baseUri);
    });
  }
}
