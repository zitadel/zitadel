import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IamComponent } from './iam.component';

describe('IamComponent', () => {
  let component: IamComponent;
  let fixture: ComponentFixture<IamComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IamComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IamComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
