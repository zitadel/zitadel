import { Injectable, Renderer2, RendererFactory2 } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';

import { KeyboardShortcutsComponent } from '../modules/keyboard-shortcuts/keyboard-shortcuts.component';

@Injectable({
  providedIn: 'root',
})
export class KeyboardShortcutsService {
  renderer: Renderer2;

  constructor(private rendererFactory2: RendererFactory2, private dialog: MatDialog) {
    this.renderer = this.rendererFactory2.createRenderer(null, null);
    console.log(this.renderer);
    this.renderer.listen('document', 'keydown', (event) => {
      const tagname = event.target.tagName;
      const exclude = ['input', 'textarea'];

      if (exclude.indexOf(tagname.toLowerCase()) === -1) {
        if (event.key === '?') {
          this.openOverviewDialog();
        }

        console.log(event);

        // SIDEWIDESHORTCUTS.f
      }
    });
  }

  public openOverviewDialog(): void {
    const dialogRef = this.dialog.open(KeyboardShortcutsComponent, {
      data: {},
      width: '600px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
      }
    });
  }
}
