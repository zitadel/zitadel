import { Component } from '@angular/core';
import { enterAnimations } from 'src/app/animations';

@Component({
  selector: 'cnsl-org-list',
  templateUrl: './org-list.component.html',
  styleUrls: ['./org-list.component.scss'],
  animations: [enterAnimations],
})
export class OrgListComponent {}
