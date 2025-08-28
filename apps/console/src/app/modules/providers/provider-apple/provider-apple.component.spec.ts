import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderAppleComponent } from './provider-apple.component';

describe('ProviderGoogleComponent', () => {
  let component: ProviderAppleComponent;
  let fixture: ComponentFixture<ProviderAppleComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderAppleComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderAppleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
