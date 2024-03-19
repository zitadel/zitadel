import {Component, Input} from '@angular/core';
import { Observable } from 'rxjs';

export interface Next {
  copyUrls: {
    label: string;
    url: string;
    downloadable?: boolean;
  }[],
  autofillLink?: string,
  configureTitle: string,
  configureDescription: string,
}

@Component({
  selector: 'cnsl-provider-next',
  templateUrl: './provider-next.component.html',
})
export class ProviderNextComponent {
  @Input({required: true}) next!: Next;
  constructor() {}
}
