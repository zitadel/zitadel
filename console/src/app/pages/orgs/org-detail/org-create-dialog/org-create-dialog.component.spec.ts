import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OrgCreateDialogComponent } from './org-create-dialog.component';

describe('OrgCreateDialogComponent', () => {
  let component: OrgCreateDialogComponent;
  let fixture: ComponentFixture<OrgCreateDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OrgCreateDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OrgCreateDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
