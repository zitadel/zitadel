import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatSortModule } from '@angular/material/sort';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

import { MembersTableComponent } from './members-table.component';

describe('MembersTableComponent', () => {
  let component: MembersTableComponent;
  let fixture: ComponentFixture<MembersTableComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [MembersTableComponent],
      imports: [NoopAnimationsModule, MatSortModule, MatTableModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MembersTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should compile', () => {
    expect(component).toBeTruthy();
  });
});
