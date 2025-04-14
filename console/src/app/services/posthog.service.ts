import { DestroyRef, Injectable, OnDestroy } from '@angular/core';
import { EnvironmentService } from './environment.service';
import posthog from 'posthog-js';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Injectable({
  providedIn: 'root',
})
export class PosthogService implements OnDestroy {
  private posthogToken?: string;
  private posthogUrl?: string;

  constructor(envService: EnvironmentService, destroyRef: DestroyRef) {
    envService.env.pipe(takeUntilDestroyed(destroyRef)).subscribe((env) => {
      this.posthogToken = env.posthog_token;
      this.posthogUrl = env.posthog_url;
      this.init();
    });
  }

  init() {
    if (this.posthogToken && this.posthogUrl) {
      posthog.init(this.posthogToken, {
        api_host: this.posthogUrl,
        session_recording: {
          maskAllInputs: true,
          maskTextSelector: '*',
        },
        enable_heatmaps: true,
        persistence: 'memory',
      });
    }
  }

  ngOnDestroy() {
    posthog.reset();
  }
}
