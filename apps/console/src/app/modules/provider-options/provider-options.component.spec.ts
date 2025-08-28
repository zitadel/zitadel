import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProviderOptionsComponent } from './provider-options.component';

describe('ProviderOptionsComponent', () => {
  let component: ProviderOptionsComponent;
  let fixture: ComponentFixture<ProviderOptionsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ProviderOptionsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ProviderOptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
