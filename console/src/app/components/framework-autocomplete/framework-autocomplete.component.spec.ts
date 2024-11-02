import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FrameworkAutocompleteComponent } from './framework-autocomplete.component';

describe('FrameworkAutocompleteComponent', () => {
  let component: FrameworkAutocompleteComponent;
  let fixture: ComponentFixture<FrameworkAutocompleteComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FrameworkAutocompleteComponent],
    });
    fixture = TestBed.createComponent(FrameworkAutocompleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
