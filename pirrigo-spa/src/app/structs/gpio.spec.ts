import { Gpio } from './gpio';

describe('Gpio', () => {
  it('should create an instance', () => {
    expect(new Gpio()).toBeTruthy();
  });
});
