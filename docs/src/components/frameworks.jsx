import React from "react";
import { Tile } from "./tile";
import frameworks from "../../frameworks.json";

export function Frameworks({}) {
  return (
    <div className="tile-wrapper">
      {frameworks.map((framework) => {
        return (
          <Tile
            title={framework.title}
            imageSource={framework.imgSrcDark}
            imageSourceLight={framework.imgSrcLight}
            link={framework.docsLink}
            external={framework.external}
          ></Tile>
        );
      })}
    </div>
  );
}
