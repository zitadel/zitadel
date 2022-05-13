import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SecretGeneratorComponent } from './secret-generator.component';

describe('OIDCConfigurationComponent', () => {
  let component: SecretGeneratorComponent;
  let fixture: ComponentFixture<SecretGeneratorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SecretGeneratorComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SecretGeneratorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
