import React from "react";
import ReactPlayer from 'react-player'

export function ResponsivePlayer({ url, ...props }) {

return (
    <div className='player-wrapper'>
        <ReactPlayer
        className='react-player'
        url={url}
        width='100%'
        height='100%'
        {...props}
        />
    </div>
    );
}