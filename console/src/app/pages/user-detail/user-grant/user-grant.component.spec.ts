import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserGrantComponent } from './user-grant.component';

describe('UserGrantComponent', () => {
  let component: UserGrantComponent;
  let fixture: ComponentFixture<UserGrantComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserGrantComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserGrantComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
