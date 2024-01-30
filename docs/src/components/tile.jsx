import React from 'react';
import styles from "../css/tile.module.css";

export function Tile({title, imageSource, link, external}) {
    return (
        <div className={styles.tile}>
            <a href={link}>
                <h4>{title}</h4>
                <img
                    className={styles.tileimg}
                    src={imageSource}
                    alt={title}
                    width={70}
                    height={70}
                />
            </a>
        </div>

    );
}
