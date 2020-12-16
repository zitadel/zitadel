import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ExternalIdpsComponent } from './external-idps.component';

describe('ExternalIdpsComponent', () => {
  let component: ExternalIdpsComponent;
  let fixture: ComponentFixture<ExternalIdpsComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ ExternalIdpsComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ExternalIdpsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
