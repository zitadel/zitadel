import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { DetailFormComponent } from './detail-form.component';

describe('DetailFormComponent', () => {
  let component: DetailFormComponent;
  let fixture: ComponentFixture<DetailFormComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [DetailFormComponent],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DetailFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
