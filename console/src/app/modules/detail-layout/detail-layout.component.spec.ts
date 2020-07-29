import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DetailLayoutComponent } from './detail-layout.component';

describe('DetailLayoutComponent', () => {
  let component: DetailLayoutComponent;
  let fixture: ComponentFixture<DetailLayoutComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DetailLayoutComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DetailLayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
