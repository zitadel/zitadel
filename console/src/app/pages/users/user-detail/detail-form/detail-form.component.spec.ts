import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { DetailFormComponent } from './detail-form.component';

describe('DetailFormComponent', () => {
  let component: DetailFormComponent;
  let fixture: ComponentFixture<DetailFormComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [DetailFormComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DetailFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should not emit submitData when form is invalid', () => {
    spyOn(component.submitData, 'emit');
    component.profileForm.controls.displayName.setValue('');
    expect(component.profileForm.invalid).toBeTrue();

    component.submitForm({} as any);
    expect(component.submitData.emit).not.toHaveBeenCalled();
  });

  it('should emit submitData when form is valid', () => {
    spyOn(component.submitData, 'emit');
    component.profileForm.controls.givenName.setValue('Given');
    component.profileForm.controls.familyName.setValue('Family');
    component.profileForm.controls.displayName.setValue('Display');
    expect(component.profileForm.valid).toBeTrue();

    component.submitForm({ id: '1' } as any);
    expect(component.submitData.emit).toHaveBeenCalledWith(jasmine.objectContaining({
      id: '1',
      displayName: 'Display',
      givenName: 'Given',
      familyName: 'Family',
    }));
  });
});
