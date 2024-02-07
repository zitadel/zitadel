import React from "react";
import styles from "../css/tile.module.css";

export function Tile({ title, imageSource, imageSourceLight, link, external }) {
  return (
    <a
      href={link}
      className={styles.tile}
      target={external ? "_blank" : "_self"}
    >
      <h4>{title}</h4>
      <img
        className={imageSourceLight ? "hideonlight" : ""}
        src={imageSource}
        alt={title}
        width={70}
        height={70}
      />
      {imageSourceLight && (
        <img
          className={imageSourceLight ? "hideondark" : ""}
          src={imageSourceLight}
          alt={title}
          width={70}
          height={70}
        />
      )}
      {external && (
        <div className={styles.external}>
          <i className="las la-external-link-alt"></i>
        </div>
      )}
    </a>
  );
}
