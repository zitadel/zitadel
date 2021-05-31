import { Component, Input, OnInit } from '@angular/core';
import { Features, FeaturesState } from 'src/app/proto/generated/zitadel/features_pb';

@Component({
  selector: 'cnsl-zitadel-tier',
  templateUrl: './zitadel-tier.component.html',
  styleUrls: ['./zitadel-tier.component.scss'],
})
export class ZitadelTierComponent implements OnInit {
  @Input() public features!: Features.AsObject;
  @Input() public iam: boolean = false;

  FeaturesState: any = FeaturesState;
  constructor() { }

  ngOnInit(): void {
  }

}
