import { VNodeTypes, ConcreteComponent } from 'vue';
import { Stubs } from './types';
export declare type CustomCreateStub = (params: {
    name: string;
    component: ConcreteComponent;
}) => ConcreteComponent;
interface StubOptions {
    name: string;
    type?: ConcreteComponent;
    renderStubDefaultSlot?: boolean;
}
export declare const registerStub: ({ source, stub, originalStub }: {
    source: ConcreteComponent;
    stub: ConcreteComponent;
    originalStub?: ConcreteComponent<{}, any, any, import("vue").ComputedOptions, import("vue").MethodOptions> | undefined;
}) => void;
export declare const getOriginalVNodeTypeFromStub: (type: ConcreteComponent) => VNodeTypes | undefined;
export declare const getOriginalStubFromSpecializedStub: (type: ConcreteComponent) => VNodeTypes | undefined;
export declare const addToDoNotStubComponents: (type: ConcreteComponent) => WeakSet<ConcreteComponent<{}, any, any, import("vue").ComputedOptions, import("vue").MethodOptions>>;
export declare const createStub: ({ name, type, renderStubDefaultSlot }: StubOptions) => import("vue").DefineComponent<any, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
    [key: string]: any;
}>, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, Record<string, any>, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<any>, {} | {
    [x: string]: any;
}>;
export declare function stubComponents(stubs?: Stubs, shallow?: boolean, renderStubDefaultSlot?: boolean): void;
export {};
