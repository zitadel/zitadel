import { ComponentPublicInstance } from 'vue';
import type { VueWrapper } from '../vueWrapper';
export declare function disableAutoUnmount(): void;
export declare function enableAutoUnmount(hook: Function): void;
export declare function trackInstance(wrapper: VueWrapper<ComponentPublicInstance>): void;
