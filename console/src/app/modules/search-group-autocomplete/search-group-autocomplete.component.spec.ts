import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { SearchGroupAutocompleteComponent } from './search-group-autocomplete.component';

describe('SearchGroupAutocompleteComponent', () => {
  let component: SearchGroupAutocompleteComponent;
  let fixture: ComponentFixture<SearchGroupAutocompleteComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [SearchGroupAutocompleteComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SearchGroupAutocompleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
