import { Component, Input } from '@angular/core';

@Component({
  selector: 'cnsl-table-actions',
  templateUrl: './table-actions.component.html',
  styleUrls: ['./table-actions.component.scss'],
  standalone: false,
})
export class TableActionsComponent {
  @Input() hasActions: boolean = false;
  constructor() {}
}
