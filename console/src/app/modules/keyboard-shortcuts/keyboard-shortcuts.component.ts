import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { Router } from '@angular/router';
import { KeyboardShortcut, ORGSHORTCUTS, SIDEWIDESHORTCUTS } from 'src/app/services/keyboard-shortcuts/keyboard-shortcuts';

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
      ['/instance', '/views', '/failed-events'].includes(this.router.url) ||
      new RegExp('/instance/policy/*').test(this.router.url)
    );
  }
}
