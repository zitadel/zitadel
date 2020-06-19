import { PasswordPatternPipe } from './password-pattern.pipe';

describe('PasswordPatternPipe', () => {
  it('create an instance', () => {
    const pipe = new PasswordPatternPipe();
    expect(pipe).toBeTruthy();
  });
});
