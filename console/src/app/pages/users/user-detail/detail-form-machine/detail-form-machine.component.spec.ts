import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { DetailFormMachineComponent } from './detail-form-machine.component';

describe('DetailFormComponent', () => {
  let component: DetailFormMachineComponent;
  let fixture: ComponentFixture<DetailFormMachineComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [DetailFormMachineComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DetailFormMachineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
