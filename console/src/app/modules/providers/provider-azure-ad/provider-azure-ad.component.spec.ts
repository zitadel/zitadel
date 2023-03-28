import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderAzureADComponent } from './provider-azure-ad.component';

describe('ProviderAzureADComponent', () => {
  let component: ProviderAzureADComponent;
  let fixture: ComponentFixture<ProviderAzureADComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderAzureADComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderAzureADComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
