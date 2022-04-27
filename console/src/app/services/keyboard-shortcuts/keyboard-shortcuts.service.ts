import { Injectable, OnDestroy, Renderer2, RendererFactory2 } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { BehaviorSubject, Observable, of, pairwise, Subject, takeUntil } from 'rxjs';

import { KeyboardShortcutsComponent } from '../../modules/keyboard-shortcuts/keyboard-shortcuts.component';
import { GrpcAuthService } from '../grpc-auth.service';
import {
  ACTIONS,
  DOMAINS,
  HOME,
  KeyboardShortcut,
  ME,
  ORG,
  PROJECTS,
  SYSTEM,
  USERGRANTS,
  USERS,
} from './keyboard-shortcuts';

@Injectable({
  providedIn: 'root',
})
export class KeyboardShortcutsService implements OnDestroy {
  private renderer: Renderer2;
  private keyPressed: BehaviorSubject<any> = new BehaviorSubject<any>(null);
  private destroy$: Subject<void> = new Subject();

  constructor(
    private rendererFactory2: RendererFactory2,
    private dialog: MatDialog,
    private router: Router,
    private authService: GrpcAuthService,
  ) {
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
          if (this.hasPermission(HOME) && this.hasFeature(HOME)) {
            this.router.navigate(HOME.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyS') {
          if (this.hasPermission(SYSTEM) && this.hasFeature(SYSTEM)) {
            this.router.navigate(SYSTEM.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyO') {
          if (this.hasPermission(ORG) && this.hasFeature(ORG)) {
            this.router.navigate(ORG.link);
          }
        }
        if (firstKey.code === 'KeyM' && secondKey.code === 'KeyE') {
          if (this.hasPermission(ME) && this.hasFeature(ME)) {
            this.router.navigate(ME.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyP') {
          if (this.hasPermission(PROJECTS) && this.hasFeature(PROJECTS)) {
            this.router.navigate(PROJECTS.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyU') {
          if (this.hasPermission(USERS) && this.hasFeature(USERS)) {
            this.router.navigate(USERS.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyA') {
          if (this.hasPermission(USERGRANTS) && this.hasFeature(USERGRANTS)) {
            this.router.navigate(USERGRANTS.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyF') {
          if (this.hasPermission(ACTIONS) && this.hasFeature(ACTIONS)) {
            this.router.navigate(ACTIONS.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyD') {
          if (this.hasPermission(DOMAINS) && this.hasFeature(DOMAINS)) {
            this.router.navigate(DOMAINS.link);
          }
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

  private hasPermission(shortcut: KeyboardShortcut): Observable<boolean> {
    if (shortcut.permissions) {
      return this.authService.isAllowed(shortcut.permissions);
    } else {
      return of(true);
    }
  }

  private hasFeature(shortcut: KeyboardShortcut): Observable<boolean> {
    if (shortcut.features) {
      return this.authService.canUseFeature(shortcut.features);
    } else {
      return of(true);
    }
  }
}
