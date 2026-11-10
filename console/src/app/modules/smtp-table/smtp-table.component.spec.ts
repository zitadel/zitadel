import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { SMTPTableComponent } from './smtp-table.component';

describe('UserTableComponent', () => {
  let component: SMTPTableComponent;
  let fixture: ComponentFixture<SMTPTableComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [SMTPTableComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SMTPTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
