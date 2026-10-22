import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription, firstValueFrom, from, timeout } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { Change } from 'src/app/proto/generated/zitadel/change_pb';

/** Same self-only history endpoint used by the native Console's MYUSER timeline. */
@Component({
  selector: 'cnsl-account-activity',
  standalone: false,
  template: `
    <h1>{{ 'ACCOUNT.ACTIVITY' | translate }}</h1>
    <p class="intro">{{ 'ACCOUNT.ACTIVITY_DESC' | translate }}</p>
    <button type="button" (click)="load(true)" [disabled]="loading">{{ 'ACTIONS.REFRESH' | translate }}</button>
    @if (failed) {
      <p role="alert">{{ 'ACCOUNT.ACTIVITY_ERROR' | translate }}</p>
    }
    @if (!loading && !failed && !changes.length) {
      <p>{{ 'ACCOUNT.ACTIVITY_EMPTY' | translate }}</p>
    }
    @for (day of activityDays; track day.key) {
      <section class="activity-day">
        <h2 class="activity-date">{{ day.date | timestampToDate | localizedDate: 'dd. MMM yyyy' }}</h2>
        @for (change of day.changes; track change.sequence) {
          <article class="activity-event">
            @if (change.editorDisplayName) {
              <cnsl-avatar
                [name]="change.editorDisplayName"
                [avatarUrl]="change.editorAvatarUrl || ''"
                [forColor]="change.editorPreferredLoginName || 'A'"
                [size]="32"
              ></cnsl-avatar>
            } @else {
              <span class="activity-fallback"></span>
            }
            <div class="activity-message">
              <p>{{ change.eventType?.localizedMessage || change.eventType?.key }}</p>
              @if (change.changeDate) {
                <time>{{ change.changeDate | timestampToDate | localizedDate: 'HH:mm' }}</time>
              }
            </div>
          </article>
        }
      </section>
    }
    @if (loading) {
      <p role="status">{{ 'ACCOUNT.LOADING' | translate }}</p>
    }
    @if (changes.length && !done) {
      <button type="button" (click)="load()" [disabled]="loading">{{ 'CHANGES.LOADMORE' | translate }}</button>
    } @else if (done && changes.length) {
      <p>{{ 'CHANGES.BOTTOM' | translate }}</p>
    }
  `,
})
export class AccountActivityComponent implements OnInit, OnDestroy {
  changes: Change.AsObject[] = [];
  loading = false;
  failed = false;
  done = false;
  private cursor = 0;
  get activityDays(): Array<{ key: string; date: { seconds: number; nanos: number }; changes: Change.AsObject[] }> {
    const days = new Map<string, { key: string; date: { seconds: number; nanos: number }; changes: Change.AsObject[] }>();
    for (const change of this.changes) {
      if (!change.changeDate) continue;
      const date = new Date(change.changeDate.seconds * 1000 + change.changeDate.nanos / 1_000_000);
      const key = `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}`;
      const day = days.get(key) ?? { key, date: change.changeDate, changes: [] };
      day.changes.push(change);
      days.set(key, day);
    }
    return [...days.values()];
  }
  private generation = 0;
  private language?: Subscription;
  constructor(
    private native: GrpcAuthService,
    private translate: TranslateService,
  ) {}
  ngOnInit(): void {
    this.language = this.translate.onLangChange.subscribe(() => void this.load(true));
    void this.load(true);
  }
  ngOnDestroy(): void {
    this.generation++;
    this.language?.unsubscribe();
  }
  async load(reset = false): Promise<void> {
    if (this.loading && !reset) return;
    const generation = ++this.generation;
    this.loading = true;
    this.failed = false;
    if (reset) {
      this.changes = [];
      this.cursor = 0;
      this.done = false;
    }
    try {
      const page = await firstValueFrom(from(this.native.listMyUserChanges(30, this.cursor)).pipe(timeout(15000)));
      if (generation !== this.generation) return;
      const last = page.resultList.at(-1)?.sequence;
      const existing = new Set(this.changes.map((change) => change.sequence));
      this.changes = [...this.changes, ...page.resultList.filter((change) => !existing.has(change.sequence))];
      this.done = page.resultList.length < 30 || last === undefined || last === this.cursor;
      if (last !== undefined) this.cursor = last;
    } catch {
      if (generation === this.generation) this.failed = true;
    } finally {
      if (generation === this.generation) this.loading = false;
    }
  }
}
