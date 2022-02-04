import { MediaMatcher } from '@angular/cdk/layout';
import { ConnectionPositionPair, Overlay, OverlayConfig, OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal } from '@angular/cdk/portal';
import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { InfoOverlayComponent } from '../modules/info-overlay/info-overlay.component';
import { StorageLocation, StorageService } from './storage.service';

interface CnslOverlay {
  id: string;
  requirements?: {
    media?: string;
    permission?: string[];
    feature?: string[];
  };
}

interface InfoOverlayConfig {}

const DEFAULT_CONFIG: InfoOverlayConfig = {
  backdropClass: 'dark-backdrop',
  panelClass: 'tm-file-preview-dialog-panel',
};

export const IntroWorkflowOverlays: CnslOverlay[] = [
  { id: 'orgswitcher', requirements: { permission: ['org.read'] } },
  { id: 'systembutton', requirements: { permission: ['iam.read'] } },
  { id: 'profilebutton' },
  { id: 'mainnav' },
];

@Injectable({
  providedIn: 'root',
})
export class OverlayService {
  public readonly currentWorkflow$: BehaviorSubject<CnslOverlay[]> = new BehaviorSubject<CnslOverlay[]>([]);
  // public readonly currentOverlayId$: BehaviorSubject<string> = new BehaviorSubject<string>('');
  public readonly nextExists$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public readonly previousExists$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private currentIndex: number | null = null;

  constructor(private mediaMatcher: MediaMatcher, private storageService: StorageService, private overlay: Overlay) {
    const media: string = '(max-width: 500px)';
    const small = this.mediaMatcher.matchMedia(media).matches;
    if (small) {
    }

    setTimeout(() => {
      const introDismissed = storageService.getItem('intro-dismissed', StorageLocation.local);
      if (!introDismissed) {
        console.log('launch intro');
        // this.currentWorkflow$.next(IntroWorkflowOverlays);
        // this.currentIndex = 0;
        // this.currentOverlayId$.next(IntroWorkflowOverlays[this.currentIndex].id);
        // this.nextExists$.next(this.currentIndex < IntroWorkflowOverlays.length - 1);
        // this.previousExists$.next(false);
        const element: HTMLElement | null = document.getElementById('orgswitchbutton');
        if (element) {
          console.log(this.overlay);

          // const overlayRef = this.overlay.create();
          // this.overlay.position().global().centerHorizontally().centerVertically();
          // .flexibleConnectedTo(element)
          // .setOrigin(element)
          // .withDefaultOffsetY(100)
          // .withDefaultOffsetX(50);
          // const infoOverlayPortal = new ComponentPortal(InfoOverlayComponent);
          // overlayRef.attach(infoOverlayPortal);

          const overlayRef = this.createOverlay(element);
          const overlayPortal = new ComponentPortal(InfoOverlayComponent);
          overlayRef.attach(overlayPortal);
        }
        // overlayRef.detach();
      }
    }, 1000);
  }

  public triggerPrevious(): void {
    if (this.currentIndex && this.currentIndex > 0) {
      this.currentIndex--;
      // this.currentOverlayId$.next(this.currentWorkflow$.value[this.currentIndex].id);
      this.nextExists$.next(this.currentIndex < this.currentWorkflow$.value.length - 1);
      this.previousExists$.next(this.currentIndex > 0);
    }
  }

  public triggerNext(): void {
    if (this.currentIndex !== null && this.currentIndex < this.currentWorkflow$.value.length) {
      this.currentIndex++;
      // this.currentOverlayId$.next(this.currentWorkflow$.value[this.currentIndex].id);
      this.nextExists$.next(this.currentIndex < this.currentWorkflow$.value.length - 1);
      this.previousExists$.next(this.currentIndex > 0);
    }
  }

  public complete(): void {
    this.currentIndex = null;
    // this.currentOverlayId$.next('');
    this.currentWorkflow$.next([]);
  }

  private createOverlay(element: HTMLElement): OverlayRef {
    // Returns an OverlayConfig
    const overlayConfig = this.getOverlayConfig(element);

    // Returns an OverlayRef
    return this.overlay.create(overlayConfig);
  }

  private getOverlayConfig(element: HTMLElement): OverlayConfig {
    // const positionStrategy = this.overlay.position().global().centerHorizontally().centerVertically();
    const positions = [
      new ConnectionPositionPair({ originX: 'start', originY: 'bottom' }, { overlayX: 'start', overlayY: 'top' }),
      new ConnectionPositionPair({ originX: 'start', originY: 'top' }, { overlayX: 'start', overlayY: 'bottom' }),
    ];
    const positionStrategy = this.overlay
      .position()
      .flexibleConnectedTo(element)
      .withPositions(positions)
      .withFlexibleDimensions(false)
      .withPush(false);

    const overlayConfig = new OverlayConfig({
      hasBackdrop: true,
      backdropClass: 'cnsl-overlay-backdrop',
      panelClass: 'cnsl-overlay-panel',
      scrollStrategy: this.overlay.scrollStrategies.block(),
      positionStrategy,
    });

    return overlayConfig;
  }
}
