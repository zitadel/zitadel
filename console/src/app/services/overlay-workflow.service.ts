import { MediaMatcher } from '@angular/cdk/layout';
import { Injectable, OnDestroy } from '@angular/core';
import { BehaviorSubject, Observable, of, Subject, switchMap, take, takeUntil, zip } from 'rxjs';

import { GrpcAuthService } from './grpc-auth.service';
import { CnslOverlayRef } from './overlay-ref';
import { OverlayService } from './overlay.service';

export interface OverlayWorkflow {
  currentIndex: number;
  currentRef?: CnslOverlayRef;
  overlays: CnslOverlay[];
}

export interface CnslOverlay {
  id: string;
  origin: string;
  toHighlight: string[];
  content: {
    i18nText: string;
  };
  requirements?: {
    media?: string;
    permission?: string[];
    feature?: string[];
  };
}

@Injectable({
  providedIn: 'root',
})
export class OverlayWorkflowService implements OnDestroy {
  public readonly currentWorkflow$: BehaviorSubject<OverlayWorkflow | null> = new BehaviorSubject<OverlayWorkflow | null>(
    null,
  );
  public destroy$: Subject<void> = new Subject();

  public openRef: CnslOverlayRef | null = null;
  public highlightedIds: { [id: string]: number } = {};
  public callback: Function | null = null;
  constructor(private mediaMatcher: MediaMatcher, overlayService: OverlayService, private authService: GrpcAuthService) {
    // const media: string = '(max-width: 500px)';
    // const small = this.mediaMatcher.matchMedia(media).matches;
    // if (small) {
    // }

    this.currentWorkflow$.pipe(takeUntil(this.destroy$)).subscribe((workflow) => {
      if (this.openRef) {
        this.openRef.close();
        this.openRef = null;
      }

      Object.keys(this.highlightedIds).forEach((id) => {
        const element = document.getElementById(id);
        if (element) {
          element.style.removeProperty('z-index');
          delete this.highlightedIds[id];
        }
      });

      const overlay = workflow?.overlays[workflow.currentIndex];
      if (overlay) {
        console.log('check req', overlay);
        this.meetsRequirements(overlay)
          .pipe(take(1))
          .subscribe((can) => {
            if (can) {
              this.openRef = overlayService.open(overlay);

              overlay.toHighlight.forEach((id) => {
                const element = document.getElementById(id);
                if (element) {
                  const oldZ = element.style.zIndex;
                  this.highlightedIds[id] = Number(oldZ);
                  element.style.zIndex = '1001';
                }
              });
            } else {
              this.nextStep();
              console.log('skip overlay');
            }
          });
      }
    });
  }

  public meetsRequirements(overlay: CnslOverlay): Observable<boolean> {
    return zip([
      this.meetsMediaRequirements(overlay),
      this.meetsPermissionRequirements(overlay),
      this.meetsFeatureRequirements(overlay),
    ] as Observable<boolean>[]).pipe(
      switchMap(([media, permission, feature]) => {
        return of(media && permission && feature);
      }),
    );
  }

  public meetsMediaRequirements(overlay: CnslOverlay): Observable<boolean> {
    const media = overlay.requirements?.media;
    if (media) {
      const matches = this.mediaMatcher.matchMedia(media).matches;
      if (!matches) {
        return of(false);
      } else {
        return of(true);
      }
    } else {
      return of(true);
    }
  }

  public meetsPermissionRequirements(overlay: CnslOverlay): Observable<boolean> {
    const regexArray = overlay.requirements?.permission;
    if (regexArray && regexArray.length) {
      return this.authService.isAllowed(regexArray);
    } else {
      return of(true);
    }
  }

  public meetsFeatureRequirements(overlay: CnslOverlay): Observable<boolean> {
    const regexArray = overlay.requirements?.feature;
    if (regexArray && regexArray.length) {
      return this.authService.canUseFeature(regexArray);
    } else {
      return of(true);
    }
  }

  public reset(): void {
    this.currentWorkflow$.next(null);
    console.log('execute cb here');
    if (this.callback) {
      this.callback();
    }
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public startWorkflow(overlays: CnslOverlay[], cb: Function | null): void {
    if (cb) {
      this.callback = cb;
    }
    this.currentWorkflow$.next({ overlays, currentIndex: 0 });
  }

  public nextStep(): void {
    const currentWorkflow = this.currentWorkflow$.value;
    if (this.nextPossible && currentWorkflow) {
      const nextIndex = currentWorkflow?.currentIndex + 1;
      this.currentWorkflow$.next({ ...currentWorkflow, currentIndex: nextIndex });
      console.log('next');
    } else {
      this.currentWorkflow$.next(null);
      console.log('finished, execute cb here');
      if (this.callback) {
        this.callback();
      }
    }
  }

  public previousStep(): void {
    const currentWorkflow = this.currentWorkflow$.value;
    if (this.previousPossible && currentWorkflow) {
      const nextIndex = currentWorkflow?.currentIndex - 1;
      this.currentWorkflow$.next({ ...currentWorkflow, currentIndex: nextIndex });
    }
  }

  public get nextPossible(): boolean {
    const currentWorkflow = this.currentWorkflow$.value;

    if (currentWorkflow) {
      const nextIndex = currentWorkflow?.currentIndex + 1;
      if (nextIndex < currentWorkflow?.overlays.length) {
        return true;
      }
    }
    return false;
  }

  public get isLast(): boolean {
    const currentWorkflow = this.currentWorkflow$.value;

    if (currentWorkflow) {
      const currentIndex = currentWorkflow?.currentIndex;

      if (currentIndex == currentWorkflow?.overlays.length - 1) {
        return true;
      }
    }
    return false;
  }

  public get previousPossible(): boolean {
    const currentWorkflow = this.currentWorkflow$.value;

    if (currentWorkflow && currentWorkflow.currentIndex > 0) {
      return true;
    } else {
      return false;
    }
  }
}
