import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderJWTCreateComponent } from './provider-jwt-create.component';

describe('ProviderJWTCreateComponent', () => {
  let component: ProviderJWTCreateComponent;
  let fixture: ComponentFixture<ProviderJWTCreateComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderJWTCreateComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderJWTCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
