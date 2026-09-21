import { ComponentFixture, fakeAsync, TestBed, tick } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap, Router } from '@angular/router';
import { of } from 'rxjs';

import { TableSearchComponent } from './table-search.component';

describe('TableSearchComponent', () => {
  let component: TableSearchComponent;
  let fixture: ComponentFixture<TableSearchComponent>;
  let router: jasmine.SpyObj<Router>;
  let emitted: Array<string>;

  function type(term: string): void {
    const input = fixture.nativeElement.querySelector('.search-input') as HTMLInputElement;
    input.value = term;
    input.dispatchEvent(new Event('input'));
  }

  /**
   * Configures the TestBed. Kept out of fakeAsync: zone.js restores the real zone once the
   * callback returns, so anything after an await would run outside the fake zone and
   * tick() would throw.
   */
  async function configure(queryParams: Record<string, string> = {}): Promise<void> {
    router = jasmine.createSpyObj<Router>('Router', ['navigate']);
    router.navigate.and.resolveTo(true);

    const route = {
      queryParamMap: of(convertToParamMap(queryParams)),
    } as unknown as ActivatedRoute;

    await TestBed.configureTestingModule({
      declarations: [TableSearchComponent],
      providers: [
        { provide: Router, useValue: router },
        { provide: ActivatedRoute, useValue: route },
      ],
    })
      .overrideComponent(TableSearchComponent, {
        set: {
          template: '<input class="search-input" [value]="value" (input)="onInput($event)" />',
        },
      })
      .compileComponents();
  }

  /** Synchronous, so it can be called from inside a fakeAsync callback. */
  function create(): void {
    fixture = TestBed.createComponent(TableSearchComponent);
    component = fixture.componentInstance;

    emitted = [];
    component.searchChanged.subscribe((term) => emitted.push(term));

    fixture.detectChanges();
  }

  describe('without a term in the url', () => {
    beforeEach(async () => {
      await configure();
    });

    it('should create', () => {
      create();
      expect(component).toBeTruthy();
    });

    it('does not emit on init', () => {
      create();
      expect(emitted).toEqual([]);
    });

    it('debounces typing into a single emission', fakeAsync(() => {
      create();

      type('m');
      type('me');
      type('mei');
      tick(299);
      expect(emitted).toEqual([]);

      tick(1);
      expect(emitted).toEqual(['mei']);
    }));

    it('trims the term before emitting', fakeAsync(() => {
      create();

      type('  meier  ');
      tick(300);

      expect(emitted).toEqual(['meier']);
    }));

    it('does not emit again when the trimmed term did not change', fakeAsync(() => {
      create();

      type('meier');
      tick(300);
      type('meier ');
      tick(300);

      expect(emitted).toEqual(['meier']);
    }));

    it('mirrors the term into the q query parameter', fakeAsync(() => {
      create();

      type('meier');
      tick(300);

      expect(router.navigate).toHaveBeenCalledWith(
        [],
        jasmine.objectContaining({ queryParams: { q: 'meier' }, queryParamsHandling: 'merge' }),
      );
    }));

    it('drops the q parameter instead of leaving it empty', fakeAsync(() => {
      create();

      type('meier');
      tick(300);
      type('');
      tick(300);

      expect(router.navigate).toHaveBeenCalledWith([], jasmine.objectContaining({ queryParams: { q: undefined } }));
    }));
  });

  describe('with a term in the url', () => {
    it('restores and emits it on init', async () => {
      await configure({ q: 'meier' });
      create();

      expect(emitted).toEqual(['meier']);
    });

    it('caps it at the proto string limit', async () => {
      await configure({ q: 'x'.repeat(250) });
      create();

      expect(emitted).toEqual(['x'.repeat(200)]);
    });
  });
});
