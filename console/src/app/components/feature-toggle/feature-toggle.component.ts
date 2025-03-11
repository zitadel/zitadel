import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { CopyToClipboardModule } from '../../directives/copy-to-clipboard/copy-to-clipboard.module';
import { CopyRowComponent } from '../copy-row/copy-row.component';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { ToggleState, ToggleStateKeys, ToggleStates } from '../features/features.component';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { Source } from 'src/app/proto/generated/zitadel/feature/v2beta/feature_pb';
import { FormsModule } from '@angular/forms';
import { GetInstanceFeaturesResponse } from 'src/app/proto/generated/zitadel/feature/v2/instance_pb';
import { FeatureFlag } from 'src/app/proto/generated/zitadel/feature/v2/feature_pb';

@Component({
  standalone: true,
  selector: 'cnsl-feature-toggle',
  templateUrl: './feature-toggle.component.html',
  styleUrls: ['./feature-toggle.component.scss'],
  imports: [
    CommonModule,
    FormsModule,
    TranslateModule,
    MatButtonModule,
    InfoSectionModule,
    MatTooltipModule,
    CopyToClipboardModule,
    CopyRowComponent,
    MatButtonToggleModule,
  ],
})
export class FeatureToggleComponent {
  @Input() featureData: Partial<GetInstanceFeaturesResponse.AsObject> = {};
  @Input() toggleStates: any = {}; // TODO: use ToggleStates, but get rid of the any warning of ngModel
  @Input() toggleStateKey: string = '';
  @Output() toggleChange = new EventEmitter<void>();

  protected ToggleState: any = ToggleState;
  protected Source: any = Source;

  get source() {
    return this.featureData[this.toggleStateKey as ToggleStateKeys]?.source;
  }

  get enabled() {
    // TODO: remove casting as not all features are a FeatureFlag
    return (this.featureData[this.toggleStateKey as ToggleStateKeys] as FeatureFlag.AsObject)?.enabled;
  }

  onToggleChange() {
    this.toggleChange.emit();
  }
}
