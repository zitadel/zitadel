import React from "react";
import { Tile } from "./tile";
import frameworks from "../../frameworks.json";

export function Frameworks({ filter }) {
	const filteredFrameworks = frameworks.filter((framework) => {
		return filter ? filter(framework) : true;
	});

	return (
		<div className="tile-wrapper">
			{filteredFrameworks.map((framework) => (
				<Tile
					key={framework.id || framework.title}
					title={framework.title}
					imageSource={framework.imgSrcDark}
					imageSourceLight={framework.imgSrcLight}
					link={framework.docsLink}
					external={framework.external}
				></Tile>
			))}
		</div>
	);
}
