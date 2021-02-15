import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RedirectUrisComponent } from './redirect-uris.component';

describe('RedirectUrisComponent', () => {
  let component: RedirectUrisComponent;
  let fixture: ComponentFixture<RedirectUrisComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RedirectUrisComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RedirectUrisComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
