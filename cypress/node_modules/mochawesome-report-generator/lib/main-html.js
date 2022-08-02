"use strict";

var escapeHtml = require('escape-html');
/**
 * Escape entities for use in HTML
 *
 * @param {string} str Input string
 *
 * @return {string}
 */


function e(str) {
  return escapeHtml(str).replace(/&#39/g, '&#x27');
}
/**
 * Render the main report HTML to a string
 *
 * @param {object} props Report properties
 * @param {string} data Raw report data
 * @param {string} inlineScripts App JS
 * @param {string} inlineStyles App CSS
 * @param {object} options App options
 * @param {string} scriptsUrl URL for app JS
 * @param {string} stylesUrl URL for app CSS
 * @param {string} title Report page title
 * @param {boolean} useInlineAssets Whether to render JS/CSS inline
 *
 * @return {string}
 */


function renderMainHTML(props) {
  var data = props.data,
      inlineScripts = props.inlineScripts,
      inlineStyles = props.inlineStyles,
      options = props.options,
      scriptsUrl = props.scriptsUrl,
      stylesUrl = props.stylesUrl,
      title = props.title,
      useInlineAssets = props.useInlineAssets;
  var styles = useInlineAssets ? "<style>".concat(inlineStyles, "</style>") : "<link rel=\"stylesheet\" href=\"".concat(stylesUrl, "\"/>");
  var scripts = useInlineAssets ? "<script type=\"text/javascript\">".concat(inlineScripts, "</script>") : "<script src=\"".concat(scriptsUrl, "\"></script>");
  var meta = '<meta charSet="utf-8"/><meta http-equiv="X-UA-Compatible" content="IE=edge"/><meta name="viewport" content="width=device-width, initial-scale=1"/>';
  var head = "<head>".concat(meta, "<title>").concat(e(title), "</title>").concat(styles, "</head>");
  var body = "<body data-raw=\"".concat(e(data), "\" data-config=\"").concat(e(JSON.stringify(options)), "\"><div id=\"report\"></div>").concat(scripts, "</body>");
  return "<html lang=\"en\">".concat(head).concat(body, "</html>");
}

module.exports = renderMainHTML;