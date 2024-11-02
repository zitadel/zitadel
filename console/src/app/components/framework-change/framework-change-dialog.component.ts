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
import { FrameworkAutocompleteComponent } from '../framework-autocomplete/framework-autocomplete.component';
import { Framework } from '../quickstart/quickstart.component';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'cnsl-framework-change-dialog',
  templateUrl: './framework-change-dialog.component.html',
  styleUrls: ['./framework-change-dialog.component.scss'],
  standalone: true,
  imports: [MatButtonModule, MatDialogModule, TranslateModule, FrameworkAutocompleteComponent],
})
export class FrameworkChangeDialogComponent {
  public framework = signal<Framework | undefined>(undefined);

  constructor(
    public dialogRef: MatDialogRef<FrameworkChangeDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.framework.set(data.framework);
  }

  public findFramework(id: string) {
    const temp = this.data.frameworks.find((f: Framework) => f.id === id);
    this.framework.set(temp);
  }

  public close() {
    this.dialogRef.close(this.framework());
  }
}
