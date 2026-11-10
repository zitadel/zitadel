import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActionsTwoAddActionTargetComponent } from './actions-two-add-action-target.component';

describe('ActionsTwoAddActionTargetComponent', () => {
  let component: ActionsTwoAddActionTargetComponent;
  let fixture: ComponentFixture<ActionsTwoAddActionTargetComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ActionsTwoAddActionTargetComponent],
    });
    fixture = TestBed.createComponent(ActionsTwoAddActionTargetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
