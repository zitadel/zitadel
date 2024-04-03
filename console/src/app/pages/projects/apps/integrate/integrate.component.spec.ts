import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { IntegrateAppComponent } from './integrate.component';

describe('IntegrateAppComponent', () => {
  let component: IntegrateAppComponent;
  let fixture: ComponentFixture<IntegrateAppComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [IntegrateAppComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IntegrateAppComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
