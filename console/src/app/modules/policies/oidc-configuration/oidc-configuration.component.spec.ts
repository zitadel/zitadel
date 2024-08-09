import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OIDCConfigurationComponent } from './oidc-configuration.component';

import { By } from '@angular/platform-browser';

describe('OIDCConfigurationComponent', () => {
  let component: OIDCConfigurationComponent;
  let fixture: ComponentFixture<OIDCConfigurationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [OIDCConfigurationComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(OIDCConfigurationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });


  it('should not allow negative access token lifetime', () => {
    const input = fixture.debugElement.query(By.css('input[name="accessTokenLifetime"]')).nativeElement;
    fixture.detectChanges();
    input.value = -1; // Attempt to set a negative value
    input.dispatchEvent(new Event('input'));



    // const fixture = TestBed.createComponent(AppComponent);
    // fixture.detectChanges();
    // const compiled = fixture.debugElement.nativeElement;
    // expect(compiled.querySelector('.content span').textContent).toContain('console app is running!');
    expect(component.form.controls['accessTokenLifetime'].value).toBeGreaterThanOrEqual(0);

    // expect(component.form.controls['accessTokenLifetime'].value).toBeGreaterThanOrEqual(0);
  });

});
