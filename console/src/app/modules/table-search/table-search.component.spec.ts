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

  async function setup(queryParams: Record<string, string> = {}): Promise<void> {
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

    fixture = TestBed.createComponent(TableSearchComponent);
    component = fixture.componentInstance;

    emitted = [];
    component.searchChanged.subscribe((term) => emitted.push(term));

    fixture.detectChanges();
  }

  afterEach(() => {
    TestBed.resetTestingModule();
  });

  it('should create', async () => {
    await setup();
    expect(component).toBeTruthy();
  });

  it('debounces typing into a single emission', fakeAsync(async () => {
    await setup();

    type('m');
    type('me');
    type('mei');
    tick(299);
    expect(emitted).toEqual([]);

    tick(1);
    expect(emitted).toEqual(['mei']);
  }));

  it('trims the term before emitting', fakeAsync(async () => {
    await setup();

    type('  meier  ');
    tick(300);

    expect(emitted).toEqual(['meier']);
  }));

  it('does not emit again when the trimmed term did not change', fakeAsync(async () => {
    await setup();

    type('meier');
    tick(300);
    type('meier ');
    tick(300);

    expect(emitted).toEqual(['meier']);
  }));

  it('mirrors the term into the q query parameter', fakeAsync(async () => {
    await setup();

    type('meier');
    tick(300);

    expect(router.navigate).toHaveBeenCalledWith(
      [],
      jasmine.objectContaining({ queryParams: { q: 'meier' }, queryParamsHandling: 'merge' }),
    );
  }));

  it('drops the q parameter instead of leaving it empty', fakeAsync(async () => {
    await setup();

    type('meier');
    tick(300);
    type('');
    tick(300);

    expect(router.navigate).toHaveBeenCalledWith([], jasmine.objectContaining({ queryParams: { q: undefined } }));
  }));

  it('restores the term from the url on init', fakeAsync(async () => {
    await setup({ q: 'meier' });

    expect(emitted).toEqual(['meier']);
  }));

  it('does not emit on init when no term is in the url', fakeAsync(async () => {
    await setup();

    expect(emitted).toEqual([]);
  }));
});
