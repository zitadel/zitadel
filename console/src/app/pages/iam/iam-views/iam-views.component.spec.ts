import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IamViewsComponent } from './iam-views.component';

describe('IamViewsComponent', () => {
  let component: IamViewsComponent;
  let fixture: ComponentFixture<IamViewsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IamViewsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IamViewsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
