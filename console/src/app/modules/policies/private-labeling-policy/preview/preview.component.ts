import { Component, Input, OnInit } from '@angular/core';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';

@Component({
  selector: 'cnsl-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
})
export class PreviewComponent implements OnInit {
  @Input() policy!: LabelPolicy.AsObject;
  @Input() label: string = 'PREVIEW';
  @Input() logoURL: string = '';
  constructor() { }

  ngOnInit(): void {
  }

}
