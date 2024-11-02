import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FrameworkChangeComponent } from './framework-change.component';

describe('FrameworkChangeComponent', () => {
  let component: FrameworkChangeComponent;
  let fixture: ComponentFixture<FrameworkChangeComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FrameworkChangeComponent],
    });
    fixture = TestBed.createComponent(FrameworkChangeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
