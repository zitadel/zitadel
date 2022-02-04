import { MediaMatcher } from '@angular/cdk/layout';
import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { StorageLocation, StorageService } from './storage.service';

interface Overlay {
  id: string;
  requirements?: {
    media?: string;
    permission?: string[];
    feature?: string[];
  };
}

export const IntroWorkflowOverlays: Overlay[] = [
  { id: 'orgswitcher', requirements: { permission: ['org.read'] } },
  { id: 'systembutton', requirements: { permission: ['iam.read'] } },
  { id: 'profilebutton' },
  { id: 'mainnav' },
];

@Injectable({
  providedIn: 'root',
})
export class OverlayService {
  public readonly currentWorkflow$: BehaviorSubject<Overlay[]> = new BehaviorSubject<Overlay[]>([]);
  public readonly currentOverlayId$: BehaviorSubject<string> = new BehaviorSubject<string>('');
  public readonly nextExists$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public readonly previousExists$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private currentIndex: number | null = null;

  constructor(private mediaMatcher: MediaMatcher, private storageService: StorageService) {
    const media: string = '(max-width: 500px)';
    const small = this.mediaMatcher.matchMedia(media).matches;
    if (small) {
    }

    setTimeout(() => {
      const introDismissed = storageService.getItem('intro-dismissed', StorageLocation.local);
      if (!introDismissed) {
        console.log('launch intro');
        this.currentWorkflow$.next(IntroWorkflowOverlays);
        this.currentIndex = 0;
        this.currentOverlayId$.next(IntroWorkflowOverlays[this.currentIndex].id);
        this.nextExists$.next(this.currentIndex < IntroWorkflowOverlays.length - 1);
        this.previousExists$.next(false);
      }
    }, 1000);
  }

  public triggerPrevious(): void {
    if (this.currentIndex && this.currentIndex > 0) {
      this.currentIndex--;
      this.currentOverlayId$.next(this.currentWorkflow$.value[this.currentIndex].id);
      this.nextExists$.next(this.currentIndex < this.currentWorkflow$.value.length - 1);
      this.previousExists$.next(this.currentIndex > 0);
    }
  }

  public triggerNext(): void {
    if (this.currentIndex !== null && this.currentIndex < this.currentWorkflow$.value.length) {
      this.currentIndex++;
      this.currentOverlayId$.next(this.currentWorkflow$.value[this.currentIndex].id);
      this.nextExists$.next(this.currentIndex < this.currentWorkflow$.value.length - 1);
      this.previousExists$.next(this.currentIndex > 0);
    }
  }

  public complete(): void {
    this.currentIndex = null;
    this.currentOverlayId$.next('');
    this.currentWorkflow$.next([]);
  }
}
