import { FormControl } from '@angular/forms';
import { emailValidator } from './validators';

describe('emailValidator', () => {
  it('should validate standard ascii email', () => {
    const control = new FormControl('test@example.com');
    expect(emailValidator(control)).toBeNull();
  });

  it('should validate email with utf8 characters', () => {
    const control = new FormControl('müller@test.com');
    expect(emailValidator(control)).toBeNull();
  });

  it('should validate email with utf8 domain', () => {
    const control = new FormControl('test@ü.com');
    expect(emailValidator(control)).toBeNull();
  });

  it('should fail for invalid email', () => {
    const control = new FormControl('invalid-email');
    expect(emailValidator(control)).not.toBeNull();
  });
});
