import { Component, Input } from '@angular/core';
import { Features, FeaturesState } from 'src/app/proto/generated/zitadel/features_pb';

@Component({
  selector: 'cnsl-zitadel-tier',
  templateUrl: './zitadel-tier.component.html',
  styleUrls: ['./zitadel-tier.component.scss'],
})
export class ZitadelTierComponent {
  @Input() public features!: Features.AsObject;
  @Input() public iam: boolean = false;

  FeaturesState: any = FeaturesState;
}
