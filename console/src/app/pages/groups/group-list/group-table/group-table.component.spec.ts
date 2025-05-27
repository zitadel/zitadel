import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { GroupTableComponent } from './group-table.component';

describe('GroupTableComponent', () => {
  let component: GroupTableComponent;
  let fixture: ComponentFixture<GroupTableComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [GroupTableComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
