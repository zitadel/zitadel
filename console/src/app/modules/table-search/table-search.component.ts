import { Component, DestroyRef, EventEmitter, inject, Input, OnInit, Output } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged, map, take } from 'rxjs/operators';

const DEBOUNCE_MS = 300;

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

  private readonly input$ = new Subject<string>();
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly destroyRef = inject(DestroyRef);

  ngOnInit(): void {
    this.route.queryParamMap
      .pipe(
        take(1),
        map((params) => params.get('q') ?? ''),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe((term) => {
        this.value = term;
        if (term) {
          this.searchChanged.emit(term);
        }
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
