import { Component, Injectable, OnInit } from '@angular/core';
import { firstValueFrom, from, timeout } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { UserGrant } from 'src/app/proto/generated/zitadel/auth_pb';
import { UserGrantState } from 'src/app/proto/generated/zitadel/user_pb';

export interface AccountApplication {
  id: string;
  name: string;
  organization: string;
  roles: string[];
  active: boolean;
}
@Injectable({ providedIn: 'root' })
export class AccountApplications {
  constructor(private native: GrpcAuthService) {}
  async list(): Promise<AccountApplication[]> {
    const grants: UserGrant.AsObject[] = [];
    let offset = 0;
    let total = Infinity;
    while (offset < total) {
      const page = await firstValueFrom(from(this.native.listMyUserGrants(100, offset, true)).pipe(timeout(15000)));
      total = page.details?.totalResult ?? Infinity;
      if (!page.resultList.length) break;
      grants.push(...page.resultList);
      offset += page.resultList.length;
      if (!page.details && page.resultList.length < 100) break;
    }
    const apps = new Map<string, AccountApplication>();
    for (const grant of grants) {
      const id = `${grant.projectId}:${grant.orgId}`;
      const app = apps.get(id) ?? {
        id,
        name: grant.projectName,
        organization: grant.orgName,
        roles: [],
        active: false,
      };
      if (grant.state === UserGrantState.USER_GRANT_STATE_ACTIVE) {
        app.active = true;
        app.roles = [...new Set([...app.roles, ...grant.roleKeysList])];
      }
      apps.set(id, app);
    }
    return [...apps.values()].sort((a, b) => a.name.localeCompare(b.name));
  }
}
@Component({
  selector: 'cnsl-account-applications',
  standalone: false,
  template: `
    <h1>{{ 'ACCOUNT.APPS' | translate }}</h1>
    <p class="intro">{{ 'ACCOUNT.APPS_DESC' | translate }}</p>
    <div class="section-heading">
      <div>
        <p class="eyebrow">{{ 'ACCOUNT.LINKS' | translate }}</p>
        <h2>{{ 'ACCOUNT.YOUR_APPS' | translate }}</h2>
      </div>
      <button type="button" (click)="load()" [disabled]="loading">{{ 'ACTIONS.REFRESH' | translate }}</button>
    </div>
    @if (loading) {
      <p role="status">{{ 'ACCOUNT.LOADING' | translate }}</p>
    } @else if (failed) {
      <p role="alert">{{ 'ACCOUNT.APPS_ERROR' | translate }}</p>
    } @else if (!apps.length) {
      <p class="empty-state">{{ 'ACCOUNT.APPS_EMPTY' | translate }}</p>
    } @else {
      <ul class="application-list">
        @for (app of apps; track app.id) {
          <li class="application-row">
            <span class="application-monogram" aria-hidden="true">{{ app.name.slice(0, 1).toUpperCase() }}</span>
            <div class="application-details">
              <h3>{{ app.name }}</h3>
              <p>{{ app.organization }}</p>
            </div>
            @if (app.roles.length) {
              <ul class="role-list">
                @for (role of app.roles; track role) {
                  <li [attr.title]="role">{{ accessLabel(role) }}</li>
                }
              </ul>
            }
            <span class="access-state">{{
              (app.active ? 'ACCOUNT.ACCESS_ACTIVE' : 'ACCOUNT.ACCESS_INACTIVE') | translate
            }}</span>
          </li>
        }
      </ul>
    }
  `,
})
export class AccountApplicationsComponent implements OnInit {
  apps: AccountApplication[] = [];
  loading = false;
  failed = false;
  constructor(private service: AccountApplications) {}
  accessLabel(role: string): string {
    const tier = /^[^.]+\.[^.]+\.([^.]+)$/.exec(role)?.[1];
    return tier ? tier[0].toUpperCase() + tier.slice(1).replace(/[_-]/g, ' ') : role;
  }
  ngOnInit(): void {
    void this.load();
  }
  async load(): Promise<void> {
    if (this.loading) return;
    this.loading = true;
    this.failed = false;
    try {
      this.apps = await this.service.list();
    } catch {
      this.failed = true;
    } finally {
      this.loading = false;
    }
  }
}
