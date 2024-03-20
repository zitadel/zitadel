import { Component, Input } from '@angular/core';
import { Observable } from 'rxjs';

export interface CopyUrl {
  label: string;
  url: string;
  downloadable?: boolean;
}

export interface Next {
  copyUrls: CopyUrl[];
  autofillLink?: string;
  configureTitle: string;
  configureDescription: string;
  configureLink: string;
}

@Component({
  selector: 'cnsl-provider-next',
  templateUrl: './provider-next.component.html',
})
export class ProviderNextComponent {
  @Input({ required: true }) next!: Next;
  constructor() {}
}
