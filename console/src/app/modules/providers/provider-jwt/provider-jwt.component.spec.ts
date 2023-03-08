import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderJWTComponent } from './provider-jwt.component';

describe('ProviderJWTComponent', () => {
  let component: ProviderJWTComponent;
  let fixture: ComponentFixture<ProviderJWTComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderJWTComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderJWTComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
