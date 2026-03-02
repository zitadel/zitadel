import { HttpClient } from '@angular/common/http';
import { JsonPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { firstValueFrom } from 'rxjs';

interface SessionState {
  authenticated: boolean;
  expiresAt?: number;
  hasRefreshToken?: boolean;
}

const BFF_BASE_URL = 'http://localhost:3001';

@Component({
  selector: 'app-root',
  imports: [JsonPipe],
  templateUrl: './app.html',
  styleUrl: './app.css',
})
export class App {
  private readonly http = inject(HttpClient);

  protected session: SessionState = { authenticated: false };
  protected userInfo: Record<string, unknown> | null = null;
  protected loading = false;
  protected error: string | null = null;

  constructor() {
    void this.refreshSession();
  }

  protected get sessionExpiresText(): string {
    if (!this.session.expiresAt) {
      return 'n/a';
    }

    return new Date(this.session.expiresAt).toLocaleString();
  }

  protected signIn(): void {
    window.location.href = `${BFF_BASE_URL}/auth/login?returnTo=/`;
  }

  protected async refreshSession(): Promise<void> {
    this.loading = true;
    this.error = null;

    try {
      const session = await firstValueFrom(
        this.http.get<SessionState>(`${BFF_BASE_URL}/session`, {
          withCredentials: true,
        }),
      );

      this.session = session ?? { authenticated: false };
      if (!this.session.authenticated) {
        this.userInfo = null;
      }
    } catch {
      this.error = 'Could not refresh session.';
    } finally {
      this.loading = false;
    }
  }

  protected async loadUserInfo(): Promise<void> {
    this.loading = true;
    this.error = null;

    try {
      const info = await firstValueFrom(
        this.http.get<Record<string, unknown>>(`${BFF_BASE_URL}/api/userinfo`, {
          withCredentials: true,
        }),
      );

      this.userInfo = info ?? null;
    } catch {
      this.error = 'Could not load userinfo through BFF.';
    } finally {
      this.loading = false;
    }
  }

  protected async logout(): Promise<void> {
    this.loading = true;
    this.error = null;

    try {
      await firstValueFrom(
        this.http.post(`${BFF_BASE_URL}/logout`, {}, { withCredentials: true }),
      );
      this.session = { authenticated: false };
      this.userInfo = null;
    } catch {
      this.error = 'Could not log out.';
    } finally {
      this.loading = false;
    }
  }
}
