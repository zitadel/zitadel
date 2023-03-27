import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatSortModule } from '@angular/material/sort';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

import { MembershipsTableComponent } from './memberships-table.component';

describe('MembershipsTableComponent', () => {
  let component: MembershipsTableComponent;
  let fixture: ComponentFixture<MembershipsTableComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [MembershipsTableComponent],
      imports: [NoopAnimationsModule, MatSortModule, MatTableModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MembershipsTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should compile', () => {
    expect(component).toBeTruthy();
  });
});
