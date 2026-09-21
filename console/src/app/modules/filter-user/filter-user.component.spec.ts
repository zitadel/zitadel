import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { ActivatedRoute, convertToParamMap, Router } from '@angular/router';
import { of } from 'rxjs';
import { SearchQuery as UserSearchQuery } from 'src/app/proto/generated/zitadel/user_pb';

import { FilterUserComponent, SubQuery } from './filter-user.component';

describe('FilterUserComponent', () => {
  let component: FilterUserComponent;
  let fixture: ComponentFixture<FilterUserComponent>;
  let emitted: Array<Array<UserSearchQuery>>;

  const checked = { checked: true } as MatCheckboxChange;

  beforeEach(async () => {
    const router = jasmine.createSpyObj<Router>('Router', ['navigate']);
    router.navigate.and.resolveTo(true);

    const route = {
      queryParamMap: of(convertToParamMap({})),
      snapshot: { queryParamMap: convertToParamMap({}) },
    } as unknown as ActivatedRoute;

    await TestBed.configureTestingModule({
      declarations: [FilterUserComponent],
      providers: [
        { provide: Router, useValue: router },
        { provide: ActivatedRoute, useValue: route },
      ],
    })
      .overrideComponent(FilterUserComponent, { set: { template: '' } })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FilterUserComponent);
    component = fixture.componentInstance;

    emitted = [];
    component.filterChanged.subscribe((queries) => emitted.push(queries as Array<UserSearchQuery>));

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  // Regression: ticking a checkbox used to send a text query with an empty value, which the
  // API rejects with "value length must be between 1 and 200 runes" and which then persisted
  // in the URL so the error came back on every reload.
  it('does not emit a text query while its input is still empty', () => {
    component.changeCheckbox(SubQuery.DISPLAYNAME, checked);
    component.emitFilter();

    expect(emitted.length).toBeGreaterThan(0);
    expect(emitted[emitted.length - 1]).toEqual([]);
  });

  it('does not emit a text query that holds only whitespace', () => {
    component.changeCheckbox(SubQuery.DISPLAYNAME, checked);
    component.setValue(SubQuery.DISPLAYNAME, component.getSubFilter(SubQuery.DISPLAYNAME), {
      target: { value: '   ' },
    });

    expect(emitted[emitted.length - 1]).toEqual([]);
  });

  it('emits the query once a value was typed', () => {
    component.changeCheckbox(SubQuery.DISPLAYNAME, checked);
    component.setValue(SubQuery.DISPLAYNAME, component.getSubFilter(SubQuery.DISPLAYNAME), {
      target: { value: 'meier' },
    });

    const last = emitted[emitted.length - 1];
    expect(last.length).toBe(1);
    expect(last[0].toObject().displayNameQuery?.displayName).toBe('meier');
  });

  it('keeps the checkbox ticked so the input stays visible', () => {
    component.changeCheckbox(SubQuery.DISPLAYNAME, checked);
    component.emitFilter();

    expect(component.getSubFilter(SubQuery.DISPLAYNAME)).toBeDefined();
  });

  it('still emits state queries, which carry no text value', () => {
    component.changeCheckbox(SubQuery.STATE, checked);
    component.emitFilter();

    const last = emitted[emitted.length - 1];
    expect(last.length).toBe(1);
    expect(last[0].toObject().stateQuery).toBeDefined();
  });

  it('emits an empty list when the filter is reset', () => {
    component.changeCheckbox(SubQuery.EMAIL, checked);
    component.setValue(SubQuery.EMAIL, component.getSubFilter(SubQuery.EMAIL), { target: { value: 'a@b.c' } });

    component.resetFilter();

    expect(emitted[emitted.length - 1]).toEqual([]);
  });
});
