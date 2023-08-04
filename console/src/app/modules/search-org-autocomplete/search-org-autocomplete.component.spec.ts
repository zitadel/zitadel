import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { SearchOrgAutocompleteComponent } from './search-org-autocomplete.component';

describe('SearchOrgComponent', () => {
  let component: SearchOrgAutocompleteComponent;
  let fixture: ComponentFixture<SearchOrgAutocompleteComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [SearchOrgAutocompleteComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SearchOrgAutocompleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
