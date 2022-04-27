import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TopViewComponent } from './top-view.component';

describe('TopViewComponent', () => {
  let component: TopViewComponent;
  let fixture: ComponentFixture<TopViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TopViewComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TopViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
