import React from 'react';

import styles from '../css/column.module.css';

export default function Column({children}) {
    return (
      <div className={styles.column}>
        {children}
      </div>
    )
}