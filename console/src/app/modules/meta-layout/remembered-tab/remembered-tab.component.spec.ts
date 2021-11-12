import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RememberedTabComponent } from './remembered-tab.component';

describe('RememberedTabComponent', () => {
  let component: RememberedTabComponent;
  let fixture: ComponentFixture<RememberedTabComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RememberedTabComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RememberedTabComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
