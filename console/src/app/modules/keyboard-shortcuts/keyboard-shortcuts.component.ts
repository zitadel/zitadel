import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { KeyboardShortcut, SIDEWIDESHORTCUTS } from 'src/app/services/keyboard-shortcuts';

@Component({
  selector: 'cnsl-keyboard-shortcuts',
  templateUrl: './keyboard-shortcuts.component.html',
  styleUrls: ['./keyboard-shortcuts.component.scss'],
})
export class KeyboardShortcutsComponent {
  public shortcuts: KeyboardShortcut[] = SIDEWIDESHORTCUTS;
  constructor(public dialogRef: MatDialogRef<KeyboardShortcutsComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
