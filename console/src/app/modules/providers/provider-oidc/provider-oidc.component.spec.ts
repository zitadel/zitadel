import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderOIDCComponent } from './provider-oidc.component';

describe('ProviderOIDCComponent', () => {
  let component: ProviderOIDCComponent;
  let fixture: ComponentFixture<ProviderOIDCComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderOIDCComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderOIDCComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
