import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderGithubComponent } from './provider-github.component';

describe('ProviderGithubComponent', () => {
  let component: ProviderGithubComponent;
  let fixture: ComponentFixture<ProviderGithubComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderGithubComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderGithubComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
