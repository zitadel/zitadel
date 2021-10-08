import { Component, Input } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'cnsl-detail-layout',
  templateUrl: './detail-layout.component.html',
  styleUrls: ['./detail-layout.component.scss'],
})
export class DetailLayoutComponent {
  @Input() backRouterLink!: RouterLink | string | any[];
  @Input() title: string | null = '';
  @Input() description: string | null = '';
  @Input() maxWidth: boolean = true;
}
