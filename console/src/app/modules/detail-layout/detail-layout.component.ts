import { Component, Input } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-detail-layout',
  templateUrl: './detail-layout.component.html',
  styleUrls: ['./detail-layout.component.scss'],
})
export class DetailLayoutComponent {
  @Input() backRouterLink!: RouterLink;
  @Input() title: string | null = '';
  @Input() description: string | null = '';
  @Input() maxWidth: boolean = true;
}
