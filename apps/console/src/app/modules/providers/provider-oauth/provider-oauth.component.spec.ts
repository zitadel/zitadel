import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderOAuthComponent } from './provider-oauth.component';

describe('ProviderOAuthComponent', () => {
  let component: ProviderOAuthComponent;
  let fixture: ComponentFixture<ProviderOAuthComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderOAuthComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderOAuthComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
