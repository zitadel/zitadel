import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap, Router } from '@angular/router';
import { of } from 'rxjs';
import { UserSuggestionService } from 'src/app/services/user-suggestion.service';

import { FilterUserGrantsComponent } from './filter-user-grants.component';

describe('FilterUserGrantsComponent', () => {
  let component: FilterUserGrantsComponent;
  let fixture: ComponentFixture<FilterUserGrantsComponent>;

  beforeEach(async () => {
    const router = jasmine.createSpyObj<Router>('Router', ['navigate']);
    router.navigate.and.resolveTo(true);

    const route = {
      queryParams: of({}),
      queryParamMap: of(convertToParamMap({})),
    } as unknown as ActivatedRoute;

    const suggestionService = jasmine.createSpyObj<UserSuggestionService>('UserSuggestionService', ['suggest']);
    suggestionService.suggest.and.resolveTo([]);

    await TestBed.configureTestingModule({
      declarations: [FilterUserGrantsComponent],
      providers: [
        { provide: Router, useValue: router },
        { provide: ActivatedRoute, useValue: route },
        { provide: UserSuggestionService, useValue: suggestionService },
      ],
    })
      .overrideComponent(FilterUserGrantsComponent, { set: { template: '' } })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FilterUserGrantsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
