import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { ToggleState, ToggleStateKeys } from '../features/features.component';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { FormsModule } from '@angular/forms';
import { Source } from '@zitadel/proto/zitadel/feature/v2/feature_pb';

@Component({
  standalone: true,
  selector: 'cnsl-feature-toggle',
  templateUrl: './feature-toggle.component.html',
  styleUrls: ['./feature-toggle.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    CommonModule,
    FormsModule,
    TranslateModule,
    MatButtonModule,
    InfoSectionModule,
    MatTooltipModule,
    CopyToClipboardModule,
    MatButtonToggleModule,
  ],
})
export class FeatureToggleComponent {
  @Input({ required: true }) toggleStateKey!: ToggleStateKeys;
  @Input({ required: true }) toggleState!: ToggleState;
  @Output() toggleChange = new EventEmitter<ToggleState>();

  protected Source = Source;

  protected get isInherited(): boolean {
    const { source } = this.toggleState;
    return source == Source.SYSTEM || source == Source.UNSPECIFIED;
  }
}
