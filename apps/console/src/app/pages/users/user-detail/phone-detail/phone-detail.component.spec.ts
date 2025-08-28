import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { PhoneDetailComponent } from './phone-detail.component';

describe('PhoneDetailComponent', () => {
  let component: PhoneDetailComponent;
  let fixture: ComponentFixture<PhoneDetailComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [PhoneDetailComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PhoneDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
