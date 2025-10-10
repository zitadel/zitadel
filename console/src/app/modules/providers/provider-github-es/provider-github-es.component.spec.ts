import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderGithubESComponent } from './provider-github-es.component';

describe('ProviderOAuthComponent', () => {
  let component: ProviderGithubESComponent;
  let fixture: ComponentFixture<ProviderGithubESComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderGithubESComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderGithubESComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
