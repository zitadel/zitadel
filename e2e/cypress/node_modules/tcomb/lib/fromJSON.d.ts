import * as t from '../.'

declare function fromJSON<T>(value: any, type: t.Type<T>): T;

export default fromJSON