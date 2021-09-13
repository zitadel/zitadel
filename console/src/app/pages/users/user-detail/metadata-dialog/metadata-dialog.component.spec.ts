import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MetadataDialogComponent } from './metadata-dialog.component';

describe('MetadataDialogComponent', () => {
  let component: MetadataDialogComponent;
  let fixture: ComponentFixture<MetadataDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [MetadataDialogComponent],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MetadataDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
