import { MediaMatcher } from '@angular/cdk/layout';
import { Overlay } from '@angular/cdk/overlay';
import { Injectable, OnDestroy } from '@angular/core';
import { BehaviorSubject, Subject, takeUntil } from 'rxjs';

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

  public openRef!: CnslOverlayRef;
  public highlightedIds: { [id: string]: number } = {};
  constructor(private mediaMatcher: MediaMatcher, private overlayService: OverlayService, private overlay: Overlay) {
    const media: string = '(max-width: 500px)';
    const small = this.mediaMatcher.matchMedia(media).matches;
    if (small) {
    }

    this.currentWorkflow$.pipe(takeUntil(this.destroy$)).subscribe((workflow) => {
      if (this.openRef) {
        this.openRef.close();
      }

      Object.keys(this.highlightedIds).forEach((id) => {
        const element = document.getElementById(id);
        if (element) {
          element.style.zIndex = `${this.highlightedIds[id]}`;
          delete this.highlightedIds[id];
        }
      });

      const overlay = workflow?.overlays[workflow.currentIndex];
      if (overlay) {
        this.openRef = overlayService.open(overlay);

        overlay.toHighlight.forEach((id) => {
          const element = document.getElementById(id);
          if (element) {
            const oldZ = element.style.zIndex;
            this.highlightedIds[id] = Number(oldZ);

            element.style.zIndex = '1001';
          }
        });
      }
    });
  }

  public reset(): void {
    this.currentWorkflow$.next(null);
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public startWorkflow(overlays: CnslOverlay[]): void {
    this.currentWorkflow$.next({ overlays, currentIndex: 0 });
  }

  public nextStep(): void {
    const currentWorkflow = this.currentWorkflow$.value;
    if (this.nextPossible && currentWorkflow) {
      const nextIndex = currentWorkflow?.currentIndex + 1;
      this.currentWorkflow$.next({ ...currentWorkflow, currentIndex: nextIndex });
    } else {
      this.currentWorkflow$.next(null);
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
      if (nextIndex <= currentWorkflow?.overlays.length) {
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
