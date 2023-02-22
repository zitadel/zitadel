import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderGoogleCreateComponent } from './provider-google-create.component';

describe('ProviderGoogleCreateComponent', () => {
  let component: ProviderGoogleCreateComponent;
  let fixture: ComponentFixture<ProviderGoogleCreateComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderGoogleCreateComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderGoogleCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
