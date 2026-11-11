import { Component, DestroyRef, EventEmitter, inject, Input, OnInit, Output } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged, filter, map } from 'rxjs/operators';

const DEBOUNCE_MS = 300;
/** Matches the max_len of the proto string queries this term is turned into. */
const MAX_TERM_LENGTH = 200;

/**
 * Free text search box for list pages. Emits a plain string and knows nothing about
 * protobuf, which lets the same component serve the v2 user list and the v1 grant list.
 * The term is mirrored into the `q` query parameter so reloads and shared links work.
 */
@Component({
  selector: 'cnsl-table-search',
  templateUrl: './table-search.component.html',
  styleUrls: ['./table-search.component.scss'],
  standalone: false,
})
export class TableSearchComponent implements OnInit {
  @Input() public placeholder = '';
  @Output() public searchChanged = new EventEmitter<string>();

  protected value = '';

  /**
   * The term currently reflected in the URL. Starts empty so an initial visit without `q`
   * does not emit, while `?q=…` still restores.
   */
  private lastSyncedTerm = '';

  private readonly input$ = new Subject<string>();
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly destroyRef = inject(DestroyRef);

  ngOnInit(): void {
    this.route.queryParamMap
      .pipe(
        // the proto string fields cap at 200, and maxlength does not cover the URL
        map((params) => (params.get('q') ?? '').slice(0, MAX_TERM_LENGTH)),
        // Skip the writes this component made itself. Without it every keystroke would
        // emit twice, and the grant list would fire a duplicate request per stroke.
        filter((term) => term !== this.lastSyncedTerm),
        distinctUntilChanged(),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe((term) => {
        this.lastSyncedTerm = term;
        this.value = term;
        this.searchChanged.emit(term);
      });

    this.input$
      .pipe(
        map((term) => term.trim()),
        debounceTime(DEBOUNCE_MS),
        distinctUntilChanged(),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe((term) => {
        this.syncToUrl(term);
        this.searchChanged.emit(term);
      });
  }

  protected onInput(event: Event): void {
    const term = (event.target as HTMLInputElement).value;
    this.value = term;
    this.input$.next(term);
  }

  protected clear(): void {
    if (!this.value) {
      return;
    }
    this.value = '';
    this.input$.next('');
  }

  private syncToUrl(term: string): void {
    this.lastSyncedTerm = term;
    this.router
      .navigate([], {
        relativeTo: this.route,
        // undefined drops the parameter instead of leaving an empty q= behind
        queryParams: { q: term || undefined },
        queryParamsHandling: 'merge',
        replaceUrl: true,
      })
      .then();
  }
}
