import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderOIDCCreateComponent } from './provider-oidc-create.component';

describe('ProviderOIDCCreateComponent', () => {
  let component: ProviderOIDCCreateComponent;
  let fixture: ComponentFixture<ProviderOIDCCreateComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderOIDCCreateComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderOIDCCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
