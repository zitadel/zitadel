import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActionsTwoAddActionConditionComponent } from './actions-two-add-action-condition.component';

describe('ActionsTwoAddActionConditionComponent', () => {
  let component: ActionsTwoAddActionConditionComponent;
  let fixture: ComponentFixture<ActionsTwoAddActionConditionComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ActionsTwoAddActionConditionComponent],
    });
    fixture = TestBed.createComponent(ActionsTwoAddActionConditionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
