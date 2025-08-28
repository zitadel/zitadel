import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderGitlabSelfHostedComponent } from './provider-gitlab-self-hosted.component';

describe('ProviderGoogleComponent', () => {
  let component: ProviderGitlabSelfHostedComponent;
  let fixture: ComponentFixture<ProviderGitlabSelfHostedComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderGitlabSelfHostedComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderGitlabSelfHostedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
