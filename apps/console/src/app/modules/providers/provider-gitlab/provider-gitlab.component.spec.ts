import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProviderGitlabComponent } from './provider-gitlab.component';

describe('ProviderGoogleComponent', () => {
  let component: ProviderGitlabComponent;
  let fixture: ComponentFixture<ProviderGitlabComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProviderGitlabComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProviderGitlabComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
