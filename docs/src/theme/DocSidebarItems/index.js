import React from 'react';
import DocSidebarItems from '@theme-original/DocSidebarItems';
import toCustomDeprecatedItemsProps from '../../utils/deprecated-items.js';

// The DocSidebarItems component is used in generated side navs for API services.
// We wrap the original to push deprecated items to the bottom and give them a CSS class.
// This lets us easily style them differently in docs/src/css/custom.css.
export default function DocSidebarItemsWrapper(props) {
    return (
        <>
            <DocSidebarItems {...toCustomDeprecatedItemsProps(props)} />
        </>
    );
}
