import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderGoogleComponent } from './provider-google.component';

describe('ProviderGoogleComponent', () => {
  let component: ProviderGoogleComponent;
  let fixture: ComponentFixture<ProviderGoogleComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderGoogleComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderGoogleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
