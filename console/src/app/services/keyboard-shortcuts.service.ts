import { Injectable, OnDestroy, Renderer2, RendererFactory2 } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { BehaviorSubject, pairwise, Subject, takeUntil } from 'rxjs';

import { KeyboardShortcutsComponent } from '../modules/keyboard-shortcuts/keyboard-shortcuts.component';
import { ACTIONS, DOMAINS, HOME, ME, ORG, PROJECTS, SYSTEM, USERGRANTS, USERS } from './keyboard-shortcuts';

@Injectable({
  providedIn: 'root',
})
export class KeyboardShortcutsService implements OnDestroy {
  private renderer: Renderer2;
  private keyPressed: BehaviorSubject<any> = new BehaviorSubject<any>(null);
  private destroy$: Subject<void> = new Subject();

  constructor(private rendererFactory2: RendererFactory2, private dialog: MatDialog, private router: Router) {
    this.renderer = this.rendererFactory2.createRenderer(null, null);
    this.renderer.listen('document', 'keydown', (event) => {
      this.keyPressed.next(event);
    });

    this.keyPressed.pipe(pairwise(), takeUntil(this.destroy$)).subscribe(([firstKey, secondKey]) => {
      const firstTagname = firstKey?.target?.tagName;
      const secondTagname = secondKey?.target?.tagName;

      const exclude = ['input', 'textarea'];

      if (
        firstKey &&
        secondKey &&
        exclude.indexOf(firstTagname?.toLowerCase()) === -1 &&
        exclude.indexOf(secondTagname?.toLowerCase()) === -1
      ) {
        if (secondKey.key === '?') {
          this.openOverviewDialog();
        }

        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyH') {
          this.router.navigate(HOME.link);
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyS') {
          this.router.navigate(SYSTEM.link);
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyO') {
          this.router.navigate(ORG.link);
        }
        if (firstKey.code === 'KeyM' && secondKey.code === 'KeyE') {
          this.router.navigate(ME.link);
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyP') {
          this.router.navigate(PROJECTS.link);
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyU') {
          this.router.navigate(USERS.link);
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyA') {
          this.router.navigate(USERGRANTS.link);
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyF') {
          this.router.navigate(ACTIONS.link);
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyD') {
          this.router.navigate(DOMAINS.link);
        }

        // if (secondKey.shiftKey && (secondKey.code === 'Digit7' || secondKey.key==='/')) {

        // }

        // SIDEWIDESHORTCUTS.f
      } else if (secondKey && exclude.indexOf(secondTagname?.toLowerCase()) === -1) {
        if (secondKey.key === '?') {
          this.openOverviewDialog();
        }
      }
    });
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public openOverviewDialog(): void {
    const dialogRef = this.dialog.open(KeyboardShortcutsComponent, {
      data: {},
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
      }
    });
  }
}
