import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QuickstartComponent } from './quickstart.component';

describe('QuickstartComponent', () => {
  let component: QuickstartComponent;
  let fixture: ComponentFixture<QuickstartComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [QuickstartComponent],
    });
    fixture = TestBed.createComponent(QuickstartComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
