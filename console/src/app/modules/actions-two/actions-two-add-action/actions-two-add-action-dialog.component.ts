import { Component, Inject, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogClose,
  MatDialogContent,
  MatDialogModule,
  MatDialogRef,
  MatDialogTitle,
} from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import { ActionsTwoAddActionTypeComponent } from './actions-two-add-action-type/actions-two-add-action-type.component';

enum Page {
  Type,
  Condition,
  Target,
}

@Component({
  selector: 'cnsl-actions-two-add-action-dialog',
  templateUrl: './actions-two-add-action-dialog.component.html',
  styleUrls: ['./actions-two-add-action-dialog.component.scss'],
  standalone: true,
  imports: [MatButtonModule, MatDialogModule, TranslateModule, ActionsTwoAddActionTypeComponent],
})
export class ActionTwoAddActionDialogComponent {
  public page = signal<Page | undefined>(Page.Type);

  constructor(
    public dialogRef: MatDialogRef<ActionTwoAddActionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  public continue() {
    const currentPage = this.page();
    if (currentPage === Page.Type) {
      this.page.set(Page.Condition);
    } else if (currentPage === Page.Condition) {
      this.page.set(Page.Target);
    } else {
      this.dialogRef.close();
    }
  }

  public close() {
    this.dialogRef.close();
  }
}
