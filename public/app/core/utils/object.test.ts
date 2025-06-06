import { sortedDeepCloneWithoutNulls } from './object';

describe('objects', () => {
  const value = {
    hello: null,
    world: {
      deeper: 10,
      foo: null,
      arr: [null, 1, 'hello'],
      value: -Infinity,
    },
    bar: undefined,
    simple: 'A',
  };

  it('returns a clean copy', () => {
    const copy = sortedDeepCloneWithoutNulls(value);
    expect(copy).toMatchObject({
      world: {
        deeper: 10,
        arr: [null, 1, 'hello'],
        value: -Infinity,
      },
      simple: 'A',
    });
    expect(value.hello).toBeNull();
    expect(value.world.foo).toBeNull();
    expect(value.bar).toBeUndefined();
  });

  it('returns a clean copy with Infinity converted to 0', () => {
    const copy = sortedDeepCloneWithoutNulls(value, true);
    expect(copy).toMatchObject({
      world: {
        deeper: 10,
        arr: [null, 1, 'hello'],
        value: 0,
      },
      simple: 'A',
    });
    expect(value.hello).toBeNull();
    expect(value.world.foo).toBeNull();
    expect(value.bar).toBeUndefined();
  });
});
