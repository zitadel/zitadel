import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActionsTwoAddActionTypeComponent } from './actions-two-add-action-type.component';

describe('ActionsTwoAddActionTypeComponent', () => {
  let component: ActionsTwoAddActionTypeComponent;
  let fixture: ComponentFixture<ActionsTwoAddActionTypeComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ActionsTwoAddActionTypeComponent],
    });
    fixture = TestBed.createComponent(ActionsTwoAddActionTypeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
