type Predicate<T> = (x: T) => boolean;
type TypeGuardPredicate<T> = (x: any) => x is T;

interface Type<T> extends Function {
  (value: T): T;
  is: TypeGuardPredicate<T>;
  displayName: string;
  meta: {
    kind: string;
    name: string;
    identity: boolean;
  };
  t: T;
}

//
// irreducible
//

interface Irreducible<T> extends Type<T> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    predicate: TypeGuardPredicate<T>;
  };
}

export function irreducible<T>(name: string, predicate: Predicate<any>): Irreducible<T>;

//
// basic types
//

export var Any: Irreducible<any>;
export var Nil: Irreducible<void | null>;
export var String: Irreducible<string>;
export var Number: Irreducible<number>;
export var Integer: Irreducible<number>;
export var Boolean: Irreducible<boolean>;
export var Array: Irreducible<Array<any>>;
export var Object: Irreducible<Object>; // FIXME restrict to POJOs
export var Function: Irreducible<Function>;
export var Error: Irreducible<Error>;
export var RegExp: Irreducible<RegExp>;
export var Date: Irreducible<Date>;

interface ApplyCommand { $apply: Function; }
interface PushCommand { $push: Array<any>; }
interface RemoveCommand { $remove: Array<string>; }
interface SetCommand { $set: any; }
interface SpliceCommand { $splice: Array<Array<any>>; }
interface SwapCommand { $swap: { from: number; to: number; }; }
interface UnshiftCommand { $unshift: Array<any>; }
interface MergeCommand { $merge: Object; }
type Command = ApplyCommand | PushCommand | RemoveCommand | SetCommand | SpliceCommand | SwapCommand | UnshiftCommand | MergeCommand;
type UpdatePatch = Command | { [key: string]: UpdatePatch };
type Update<T> = (instance: T, spec: UpdatePatch) => T;

type Constructor<T> = Type<T> | Function;

//
// refinement
//

interface Refinement<T> extends Type<T> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    type: Constructor<T>;
    predicate: TypeGuardPredicate<T>;
  };
  update: Update<T>;
}

export function refinement<T>(type: Constructor<T>, predicate: Predicate<T>, name?: string): Refinement<T>;

//
// struct
//

type StructProps = { [key: string]: Constructor<any> };
type StructMixin = StructProps | Struct<any> | Interface<any>;

interface Struct<T> extends Type<T> {
  new(value: T): T;
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    props: StructProps;
    strict: boolean;
    defaultProps: object;
  };
  update: Update<T>;
  extend<E extends T>(mixins: StructMixin | Array<StructMixin>, name?: string | StructOptions): Struct<E>;
}

type StructOptions = {
  name?: string,
  strict?: boolean,
  defaultProps?: object
};

export function struct<T>(props: StructProps, name?: string | StructOptions): Struct<T>;

//
// interface
//

interface Interface<T> extends Type<T> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    props: StructProps;
    strict: boolean;
  };
  update: Update<T>;
  extend<E extends T>(mixins: StructMixin | Array<StructMixin>, name?: string | StructOptions): Struct<E>;
}

export function interface<T>(props: StructProps, name?: string | StructOptions): Interface<T>;

//
// list
//

interface List<T> extends Type<Array<T>> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    type: Constructor<T>;
  };
  update: Update<Array<T>>;
}

export function list<T>(type: Constructor<T>, name?: string): List<T>;

//
// dict combinator
//

interface Dict<T> extends Type<{ [key: string]: T; }> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    domain: Constructor<string>;
    codomain: T;
  };
  update: Update<{ [key: string]: T; }>;
}

export function dict<T>(domain: Constructor<string>, codomain: Constructor<T>, name?: string): Dict<T>;

//
// enums combinator
//

interface Enums extends Type<string> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    map: Object;
  };
}

interface EnumsFunction extends Function {
  (map: Object, name?: string): Enums;
  of(enums: string, name?: string): Enums;
  of(enums: Array<string>, name?: string): Enums;
}

export var enums: EnumsFunction;

//
// maybe combinator
//

interface Maybe<T> extends Type<void | T> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    type: Constructor<T>;
  };
  update: Update<void | T>;
}

export function maybe<T>(type: Constructor<T>, name?: string): Maybe<T>;

//
// tuple combinator
//

interface Tuple<T> extends Type<T> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    types: Array<Constructor<any>>;
  };
  update: Update<T>;
}

export function tuple<T>(types: Array<Constructor<any>>, name?: string): Tuple<T>;

//
// union combinator
//

interface Union<T> extends Type<T> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    types: Array<Constructor<T>>;
  };
  update: Update<T>;
  dispatch(x: any): Constructor<T>;
}

export function union<T>(types: Array<Constructor<T>>, name?: string): Union<T>;

//
// intersection combinator
//

interface Intersection<T> extends Type<T> {
  meta: {
    kind: string;
    name: string;
    identity: boolean;
    types: Array<Constructor<any>>;
  };
  update: Update<T>;
}

export function intersection<T>(types: Array<Constructor<any>>, name?: string): Intersection<T>;

//
// declare combinator
//

interface Declare<T> extends Type<T> {
  update: Update<T>;
  define(type: Type<any>): void;
}

export function declare<T>(name?: string): Declare<T>;

//
// other exports
//

export function is<T>(x: any, type: Constructor<T>): boolean;
type LazyMessage = () => string;
export function assert(guard: boolean, message?: string | LazyMessage): void;
export function fail(message: string): void;
export function isType<T>(x: Constructor<T>): boolean;
export function getTypeName<T>(x: Constructor<T>): string;
export function mixin<T, S>(target: T, source: S, overwrite?: boolean): T & S;
type Function1 = (x: any) => any;
type Clause = Constructor<any> | Function1;
export function match(x: any, ...clauses: Array<Clause>): any; // FIXME
export var update: Update<Object>;
