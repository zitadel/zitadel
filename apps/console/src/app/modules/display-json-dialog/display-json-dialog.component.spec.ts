import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DisplayJsonDialogComponent } from './display-json-dialog.component';

describe('DisplayJsonDialogComponent', () => {
  let component: DisplayJsonDialogComponent;
  let fixture: ComponentFixture<DisplayJsonDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [DisplayJsonDialogComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(DisplayJsonDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
