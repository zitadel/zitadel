import { AsyncPipe, NgIf, UpperCasePipe } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { MatTooltipModule } from '@angular/material/tooltip';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { ToggleStateKeys, ToggleStates } from '../features/features.component';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { FormsModule } from '@angular/forms';
import { Source } from '@zitadel/proto/zitadel/feature/v2/feature_pb';
import { ReplaySubject } from 'rxjs';
import { map } from 'rxjs/operators';

@Component({
  standalone: true,
  selector: 'cnsl-feature-toggle',
  templateUrl: './feature-toggle.component.html',
  styleUrls: ['./feature-toggle.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    MatButtonToggleModule,
    UpperCasePipe,
    TranslateModule,
    FormsModule,
    MatTooltipModule,
    InfoSectionModule,
    AsyncPipe,
    NgIf,
  ],
})
export class FeatureToggleComponent<TKey extends ToggleStateKeys, TValue extends ToggleStates[TKey]> {
  @Input({ required: true }) toggleStateKey!: TKey;
  @Input({ required: true })
  set toggleState(toggleState: TValue) {
    // we copy the toggleState so we can mutate it
    this.toggleState$.next(structuredClone(toggleState));
  }

  @Output() readonly toggleChange = new EventEmitter<TValue>();

  protected readonly Source = Source;
  protected readonly toggleState$ = new ReplaySubject<TValue>(1);
  protected readonly isInherited$ = this.toggleState$.pipe(
    map(({ source }) => source == Source.SYSTEM || source == Source.UNSPECIFIED),
  );
}
