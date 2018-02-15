const html = require("choo/html");
const style = require("./style");

module.exports = function(content) {
  return html`
    <body class="${style}">
      <div class="content">
        ${content}
      </div>
    </body>
  `;
};
