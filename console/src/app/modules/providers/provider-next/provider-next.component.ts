import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';

export interface CopyUrl {
  label: string;
  url: string;
  downloadable?: boolean;
}

@Component({
  selector: 'cnsl-provider-next',
  templateUrl: './provider-next.component.html',
  styleUrls: ['./provider-next.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ProviderNextComponent {
  @Input() copyUrls: CopyUrl[] | null = null;
  @Input() autofillLink: string | null = null;
  @Input({ required: true }) activateLink!: string | null;
  @Input({ required: true }) configureProvider!: boolean;
  @Input({ required: true }) configureTitle!: string;
  @Input({ required: true }) configureDescription!: string;
  @Input() configureLink?: string;
  @Input({ required: true }) expanded!: boolean;
  @Output() activate = new EventEmitter<void>();
}
