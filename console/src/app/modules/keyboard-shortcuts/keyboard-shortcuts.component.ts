import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { KeyboardShortcut, ORGSHORTCUTS, SIDEWIDESHORTCUTS } from 'src/app/services/keyboard-shortcuts';

@Component({
  selector: 'cnsl-keyboard-shortcuts',
  templateUrl: './keyboard-shortcuts.component.html',
  styleUrls: ['./keyboard-shortcuts.component.scss'],
})
export class KeyboardShortcutsComponent {
  public orgShortcuts: KeyboardShortcut[] = ORGSHORTCUTS;
  public shortcuts: KeyboardShortcut[] = SIDEWIDESHORTCUTS;
  constructor(
    public dialogRef: MatDialogRef<KeyboardShortcutsComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private router: Router,
  ) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public get isNotOnSystem(): boolean {
    return !(
      ['/system', '/views', '/failed-events'].includes(this.router.url) ||
      new RegExp('/system/policy/*').test(this.router.url)
    );
  }
}
