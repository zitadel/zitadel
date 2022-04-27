import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InfoOverlayComponent } from './info-overlay.component';

describe('InfoOverlayComponent', () => {
  let component: InfoOverlayComponent;
  let fixture: ComponentFixture<InfoOverlayComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ InfoOverlayComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(InfoOverlayComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
