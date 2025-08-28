import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PersonalAccessTokensComponent } from './personal-access-tokens.component';

describe('PersonalAccessTokensComponent', () => {
  let component: PersonalAccessTokensComponent;
  let fixture: ComponentFixture<PersonalAccessTokensComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [PersonalAccessTokensComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PersonalAccessTokensComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
