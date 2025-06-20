import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { GroupGrantCreateComponent } from './group-grant-create.component';

describe('GroupGrantCreateComponent', () => {
  let component: GroupGrantCreateComponent;
  let fixture: ComponentFixture<GroupGrantCreateComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [GroupGrantCreateComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupGrantCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
