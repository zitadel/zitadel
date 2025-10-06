import React, { useState } from "react";
import OriginalApiItem from "@theme-original/ApiItem";

export default function ApiItemWrapper(props) {
  const [feedback, setFeedback] = useState(null);
  const [comment, setComment] = useState("");

  const sendPlausibleEvent = (feedbackValue, commentText = "") => {
    if (typeof window.plausible === "function") {
      window.plausible("Feedback Submitted", {
        props: {
          feedback: feedbackValue,
          comment: commentText,
          page: window.location.pathname,
        },
      });
    }
  };

  const handleYes = () => {
    setFeedback("yes");
    sendPlausibleEvent("yes");
  };

  const handleNo = () => {
    setFeedback("no");
    sendPlausibleEvent("no");
  };

  const handleSubmitComment = () => {
    sendPlausibleEvent("no", comment);
    setFeedback("submitted");
  };

  return (
    <>
      <OriginalApiItem {...props} />

      <div className="mt-8 pt-4 border-t border-card-border">
        <p className="mb-2 font-medium text-lg">Was this page useful?</p>

        {feedback === null && (
          <div className="flex gap-4">
            <button
              className="button button--sm button--primary"
              onClick={handleYes}
            >
              Yes
            </button>
            <button
              className="button button--sm button--secondary"
              onClick={handleNo}
            >
              No
            </button>
          </div>
        )}

        {feedback === "no" && (
          <div className="mt-2 flex flex-col gap-2">
            <textarea
              className="textarea textarea-bordered w-full max-w-[50%] resize-none"
              rows="3"
              placeholder="How can we improve this page?"
              value={comment}
              onChange={(e) => setComment(e.target.value)}
            />
            <div className="flex gap-2">
              <button
                className="button button--sm button--primary"
                onClick={handleSubmitComment}
              >
                Submit
              </button>
              <button
                className="button button--sm button--secondary"
                onClick={() => {
                  setFeedback(null);
                  setComment("");
                }}
              >
                Go Back
              </button>
            </div>
          </div>
        )}

        {feedback === "yes" && <p className="mt-2">Thanks for your feedback!</p>}
        {feedback === "submitted" && (
          <p className="mt-2">Thanks for your feedback!</p>
        )}
      </div>
    </>
  );
}
