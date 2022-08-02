import { FunctionalComponent, ComponentPublicInstance, ComponentOptionsWithObjectProps, ComponentOptionsWithArrayProps, ComponentOptionsWithoutProps, ExtractPropTypes, VNodeProps, ComponentOptionsMixin, DefineComponent, MethodOptions, AllowedComponentProps, ComponentCustomProps, ExtractDefaultPropTypes, EmitsOptions, ComputedOptions, ComponentPropsOptions, Prop } from 'vue';
import { MountingOptions } from './types';
import { VueWrapper } from './vueWrapper';
declare type PublicProps = VNodeProps & AllowedComponentProps & ComponentCustomProps;
declare type ComponentMountingOptions<T> = T extends DefineComponent<infer PropsOrPropOptions, any, infer D, any, any> ? MountingOptions<Partial<ExtractDefaultPropTypes<PropsOrPropOptions>> & Omit<Readonly<ExtractPropTypes<PropsOrPropOptions>> & PublicProps, keyof ExtractDefaultPropTypes<PropsOrPropOptions>>, D> & Record<string, any> : MountingOptions<any>;
export declare function mount<V>(originalComponent: {
    new (...args: any[]): V;
    __vccOpts: any;
}, options?: MountingOptions<any> & Record<string, any>): VueWrapper<ComponentPublicInstance<V>>;
export declare function mount<V, P>(originalComponent: {
    new (...args: any[]): V;
    __vccOpts: any;
    defaultProps?: Record<string, Prop<any>> | string[];
}, options?: MountingOptions<P & PublicProps> & Record<string, any>): VueWrapper<ComponentPublicInstance<V>>;
export declare function mount<V>(originalComponent: {
    new (...args: any[]): V;
    registerHooks(keys: string[]): void;
}, options?: MountingOptions<any> & Record<string, any>): VueWrapper<ComponentPublicInstance<V>>;
export declare function mount<V, P>(originalComponent: {
    new (...args: any[]): V;
    props(Props: P): any;
    registerHooks(keys: string[]): void;
}, options?: MountingOptions<P & PublicProps> & Record<string, any>): VueWrapper<ComponentPublicInstance<V>>;
export declare function mount<Props, E extends EmitsOptions = {}>(originalComponent: FunctionalComponent<Props, E>, options?: MountingOptions<Props & PublicProps> & Record<string, any>): VueWrapper<ComponentPublicInstance<Props>>;
export declare function mount<PropsOrPropOptions = {}, RawBindings = {}, D = {}, C extends ComputedOptions = ComputedOptions, M extends MethodOptions = MethodOptions, Mixin extends ComponentOptionsMixin = ComponentOptionsMixin, Extends extends ComponentOptionsMixin = ComponentOptionsMixin, E extends EmitsOptions = Record<string, any>, EE extends string = string, PP = PublicProps, Props = Readonly<ExtractPropTypes<PropsOrPropOptions>>, Defaults = ExtractDefaultPropTypes<PropsOrPropOptions>>(component: DefineComponent<PropsOrPropOptions, RawBindings, D, C, M, Mixin, Extends, E, EE, PP, Props, Defaults>, options?: MountingOptions<Partial<Defaults> & Omit<Props & PublicProps, keyof Defaults>, D> & Record<string, any>): VueWrapper<InstanceType<DefineComponent<PropsOrPropOptions, RawBindings, D, C, M, Mixin, Extends, E, EE, PP, Props, Defaults>>>;
export declare function mount<T extends DefineComponent<any, any, any, any>>(component: T, options?: ComponentMountingOptions<T>): VueWrapper<InstanceType<T>>;
export declare function mount<Props = {}, RawBindings = {}, D = {}, C extends ComputedOptions = {}, M extends Record<string, Function> = {}, E extends EmitsOptions = Record<string, any>, Mixin extends ComponentOptionsMixin = ComponentOptionsMixin, Extends extends ComponentOptionsMixin = ComponentOptionsMixin, EE extends string = string>(componentOptions: ComponentOptionsWithoutProps<Props, RawBindings, D, C, M, E, Mixin, Extends, EE>, options?: MountingOptions<Props & PublicProps, D>): VueWrapper<ComponentPublicInstance<Props, RawBindings, D, C, M, E, VNodeProps & Props>> & Record<string, any>;
export declare function mount<PropNames extends string, RawBindings, D, C extends ComputedOptions = {}, M extends Record<string, Function> = {}, E extends EmitsOptions = Record<string, any>, Mixin extends ComponentOptionsMixin = ComponentOptionsMixin, Extends extends ComponentOptionsMixin = ComponentOptionsMixin, EE extends string = string, Props extends Readonly<{
    [key in PropNames]?: any;
}> = Readonly<{
    [key in PropNames]?: any;
}>>(componentOptions: ComponentOptionsWithArrayProps<PropNames, RawBindings, D, C, M, E, Mixin, Extends, EE, Props>, options?: MountingOptions<Props & PublicProps, D>): VueWrapper<ComponentPublicInstance<Props, RawBindings, D, C, M, E>>;
export declare function mount<PropsOptions extends Readonly<ComponentPropsOptions>, RawBindings, D, C extends ComputedOptions = {}, M extends Record<string, Function> = {}, E extends EmitsOptions = Record<string, any>, Mixin extends ComponentOptionsMixin = ComponentOptionsMixin, Extends extends ComponentOptionsMixin = ComponentOptionsMixin, EE extends string = string>(componentOptions: ComponentOptionsWithObjectProps<PropsOptions, RawBindings, D, C, M, E, Mixin, Extends, EE>, options?: MountingOptions<ExtractPropTypes<PropsOptions> & PublicProps, D>): VueWrapper<ComponentPublicInstance<ExtractPropTypes<PropsOptions>, RawBindings, D, C, M, E, VNodeProps & ExtractPropTypes<PropsOptions>>>;
export declare const shallowMount: typeof mount;
export {};
