import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppDetailComponent } from './app-detail.component';

describe('AppDetailComponent', () => {
  let component: AppDetailComponent;
  let fixture: ComponentFixture<AppDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [AppDetailComponent],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
