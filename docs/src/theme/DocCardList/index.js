import React from 'react';
import DocCardList from '@theme-original/DocCardList';
import toCustomDeprecatedItemsProps from "../../utils/deprecated-items";

// The DocCardList component is used in generated index pages for API services.
// We customize it for deprecated items.
export default function DocCardListWrapper(props) {
    return (
        <>
            <DocCardList {...toCustomDeprecatedItemsProps(props)} />
        </>
    );
}
