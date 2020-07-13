import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppCreateComponent } from './app-create.component';

describe('AppCreateComponent', () => {
  let component: AppCreateComponent;
  let fixture: ComponentFixture<AppCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [AppCreateComponent],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
