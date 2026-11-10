import React from 'react';

export function B2B() {
  return (
    <div className="flexrowbetween">
      <div className="b2borg">
        <span>Octagon <small>(owner)</small></span>
        <div className="b2bproject">
          <span>Portal Project</span>

          <div className="b2bapp">
            <strong>WEBAPP</strong>
          </div>

          <span className="b2bprojectrole">reader, writer, admin</span>
        </div>

        <div className="b2buser">
          Bill <small>(admin)</small>
        </div>
      </div>

      <svg className="arrowright" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
      
      <div className="b2borg">
        <span>Pentagon</span>
        <div className="b2bproject">
          <span>Portal Project <small>(grant)</small></span>

          <div className="b2bapp">
            <strong>WEBAPP</strong>
          </div>

          <span className="b2bprojectrole">reader, writer</span>
        </div>

        <div className="b2buser">
          Dimitri <small>(writer)</small>
        </div>

        <div className="b2buser">
          Michael <small>(reader)</small>
        </div>
      </div>
    </div>
  );
}
