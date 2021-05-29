import { Component, Input, OnInit } from '@angular/core';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';

import { Theme } from '../private-labeling-policy.component';

@Component({
  selector: 'cnsl-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
})
export class PreviewComponent implements OnInit {
  @Input() policy!: LabelPolicy.AsObject;
  @Input() label: string = 'PREVIEW';
  @Input() logoURL: string = '';
  @Input() theme: Theme = Theme.DARK;
  Theme: any = Theme;
  constructor() { }

  ngOnInit(): void {
  }
}
