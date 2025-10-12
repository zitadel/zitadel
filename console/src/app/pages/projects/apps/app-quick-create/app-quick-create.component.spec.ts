import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AppQuickCreateComponent } from './app-quick-create.component';

describe('AppQuickCreateComponent', () => {
  let component: AppQuickCreateComponent;
  let fixture: ComponentFixture<AppQuickCreateComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [AppQuickCreateComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppQuickCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
