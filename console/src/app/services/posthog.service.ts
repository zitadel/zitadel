import { Injectable } from '@angular/core';
import { EnvironmentService } from './environment.service';
import posthog from 'posthog-js';

@Injectable({
  providedIn: 'root',
})
export class PosthogService {
  private posthogToken?: string;
  private posthogUrl?: string;

  constructor(private envService: EnvironmentService) {
    this.envService.env.subscribe((env) => {
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
}
