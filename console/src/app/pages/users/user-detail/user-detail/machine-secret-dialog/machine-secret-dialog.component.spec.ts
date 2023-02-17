import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MachineSecretDialogComponent } from './machine-secret-dialog.component';

describe('MachineSecretDialogComponent', () => {
  let component: MachineSecretDialogComponent;
  let fixture: ComponentFixture<MachineSecretDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [MachineSecretDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MachineSecretDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
