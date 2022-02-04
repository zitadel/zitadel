import { MediaMatcher } from '@angular/cdk/layout';
import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { StorageLocation, StorageService } from './storage.service';

interface Overlay {
  id: string;
  workflowId: string;
  requirements?: {
    media?: string;
    permission?: string[];
    feature?: string[];
  };
}

export const IntroWorkflowOverlays: Overlay[] = [
  { id: 'orgswitcher', workflowId: 'intro', requirements: { permission: ['org.read'] } },
  { id: 'systembutton', workflowId: 'intro', requirements: { permission: ['iam.read'] } },
  { id: 'profile', workflowId: 'intro' },
];

@Injectable({
  providedIn: 'root',
})
export class OverlayService {
  public readonly currentWorkflow$: BehaviorSubject<Overlay[]> = new BehaviorSubject<Overlay[]>([]);
  public readonly currentOverlayId$: BehaviorSubject<string> = new BehaviorSubject<string>('');

  private currentIndex: number | null = null;

  constructor(private mediaMatcher: MediaMatcher, private storageService: StorageService) {
    const media: string = '(max-width: 500px)';
    const small = this.mediaMatcher.matchMedia(media).matches;
    if (small) {
    }

    const introDismissed = storageService.getItem('intro-dismissed', StorageLocation.local);
    if (!introDismissed) {
      this.currentWorkflow$.next(IntroWorkflowOverlays);
      this.currentIndex = 0;
      this.currentOverlayId$.next(IntroWorkflowOverlays[this.currentIndex].id);
    }
  }

  public triggerNext(): void {
    // this.currentIndex++;
  }
}
