import React from "react";

// This function changes a ListComponents input properties.
// Deprecated items are pushed to the bottom of the list and its labels are given the CSS class zitadel-lifecycle-deprecated.
// They are styled in docs/src/css/custom.css.
export default function (props) {
    const { items = [], ...rest } = props;
    if (!Array.isArray(items)) {
        // Do nothing if items is not an array
        return props;
    }
    const withDeprecated = [...items].map(({className, label, ...itemRest}) => {
        const zitadelLifecycleDeprecated = className?.indexOf('menu__list-item--deprecated') > -1
        const wrappedLabel = <span className={zitadelLifecycleDeprecated ? "zitadel-lifecycle-deprecated" : undefined}>{label}</span>
        return {
            zitadelLifecycleDeprecated: zitadelLifecycleDeprecated,
            ...itemRest,
            className,
            label: wrappedLabel,
        };
    });
    const sortedItems = [...withDeprecated].sort((a, b) => {
        return a.zitadelLifecycleDeprecated - b.zitadelLifecycleDeprecated;
    });
    return {
        ...rest,
        items: sortedItems,
    };
}
