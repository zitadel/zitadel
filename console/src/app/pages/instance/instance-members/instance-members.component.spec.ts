import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatSortModule } from '@angular/material/sort';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

import { InstanceMembersComponent } from './instance-members.component';

describe('IamMembersComponent', () => {
  let component: InstanceMembersComponent;
  let fixture: ComponentFixture<InstanceMembersComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [InstanceMembersComponent],
      imports: [NoopAnimationsModule, MatSortModule, MatTableModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InstanceMembersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should compile', () => {
    expect(component).toBeTruthy();
  });
});
