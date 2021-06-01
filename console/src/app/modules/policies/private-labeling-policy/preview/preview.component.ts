import { Component, Input, OnInit, SimpleChanges } from '@angular/core';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';

import { Preview, Theme } from '../private-labeling-policy.component';

@Component({
  selector: 'cnsl-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
})
export class PreviewComponent implements OnInit {
  @Input() preview: Preview = Preview.PREVIEW;
  @Input() policy!: LabelPolicy.AsObject;
  @Input() label: string = 'PREVIEW';
  @Input() images: { [imagekey: string]: any; } = {};
  @Input() theme: Theme = Theme.DARK;
  public Theme: any = Theme;
  public Preview: any = Preview;
  constructor() { }

  ngOnInit(): void {
    console.log(this.images);
  }

  ngOnChanges(changes: SimpleChanges): void {
    //Called before any other lifecycle hook. Use it to inject dependencies, but avoid any serious work here.
    //Add '${implements OnChanges}' to the class.
    console.log(this.images);

  }
}
