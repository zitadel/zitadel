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
  INSTANCE,
  KeyboardShortcut,
  ME,
  ORG,
  ORGSETTINGS,
  PROJECTS,
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
          if (this.hasPermission(HOME)) {
            this.router.navigate(HOME.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyI') {
          if (this.hasPermission(INSTANCE)) {
            this.router.navigate(INSTANCE.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyO') {
          if (this.hasPermission(ORG)) {
            this.router.navigate(ORG.link);
          }
        }
        if (firstKey.code === 'KeyM' && secondKey.code === 'KeyE') {
          if (this.hasPermission(ME)) {
            this.router.navigate(ME.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyP') {
          if (this.hasPermission(PROJECTS)) {
            this.router.navigate(PROJECTS.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyU') {
          if (this.hasPermission(USERS)) {
            this.router.navigate(USERS.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyA') {
          if (this.hasPermission(USERGRANTS)) {
            this.router.navigate(USERGRANTS.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyF') {
          if (this.hasPermission(ACTIONS)) {
            this.router.navigate(ACTIONS.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyD') {
          if (this.hasPermission(DOMAINS)) {
            this.router.navigate(DOMAINS.link);
          }
        }
        if (firstKey.code === 'KeyG' && secondKey.code === 'KeyS') {
          if (this.hasPermission(ORGSETTINGS)) {
            this.router.navigate(ORGSETTINGS.link);
          }
        }
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
    this.dialog.open(KeyboardShortcutsComponent, {
      width: '400px',
    });
  }

  private hasPermission(shortcut: KeyboardShortcut): Observable<boolean> {
    if (shortcut.permissions) {
      return this.authService.isAllowed(shortcut.permissions);
    } else {
      return of(true);
    }
  }
}
