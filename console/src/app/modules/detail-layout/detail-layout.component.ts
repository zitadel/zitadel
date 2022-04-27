import { Component, Input } from '@angular/core';

@Component({
  selector: 'cnsl-detail-layout',
  templateUrl: './detail-layout.component.html',
  styleUrls: ['./detail-layout.component.scss'],
})
export class DetailLayoutComponent {
  @Input() hasBackButton: boolean = true;
  @Input() title: string | null = '';
  @Input() description: string | null = '';
}
