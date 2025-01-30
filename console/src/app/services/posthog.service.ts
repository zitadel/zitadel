import { Injectable, OnDestroy } from '@angular/core';
import { EnvironmentService } from './environment.service';
import { Subscription } from 'rxjs';
import posthog from 'posthog-js';

@Injectable({
  providedIn: 'root',
})
export class PosthogService implements OnDestroy {
  private posthogToken?: string;
  private posthogUrl?: string;
  private envSubscription: Subscription;

  constructor(private envService: EnvironmentService) {
    this.envSubscription = this.envService.env.subscribe((env) => {
      this.posthogToken = env.posthog_token;
      this.posthogUrl = env.posthog_url;
      this.init();
    });
  }

  async init() {
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
    if (this.envSubscription) {
      this.envSubscription.unsubscribe();
    }

    posthog.reset();
  }
}
