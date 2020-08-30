import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MachineKeysComponent } from './machine-keys.component';

describe('MachineKeysComponent', () => {
  let component: MachineKeysComponent;
  let fixture: ComponentFixture<MachineKeysComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MachineKeysComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MachineKeysComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
