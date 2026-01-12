import { Component, inject, linkedSignal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { FrameworkAutocompleteComponent } from '../framework-autocomplete/framework-autocomplete.component';
import { TranslateModule } from '@ngx-translate/core';
import { frameworks } from 'src/app/utils/framework';

export type FrameworkChangeDialogData = (typeof frameworks)[number] | null;
export type FrameworkChangeDialogResult = FrameworkChangeDialogData;

@Component({
  selector: 'cnsl-framework-change-dialog',
  templateUrl: './framework-change-dialog.component.html',
  styleUrls: ['./framework-change-dialog.component.scss'],
  imports: [MatButtonModule, MatDialogModule, TranslateModule, FrameworkAutocompleteComponent],
})
export class FrameworkChangeDialogComponent {
  protected readonly frameworks = frameworks;
  private readonly dialogRef = inject(MatDialogRef<FrameworkChangeDialogComponent>);
  private readonly data = inject<FrameworkChangeDialogData>(MAT_DIALOG_DATA);

  protected framework = linkedSignal(() => this.data);

  public findFramework(frameworkId: string) {
    const framework = frameworks.find((f) => f.id === frameworkId);
    this.framework.set(framework ?? null);
  }

  public close() {
    this.dialogRef.close(this.framework());
  }
}
