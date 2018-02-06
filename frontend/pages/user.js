const html = require("choo/html");

const style = require("../lib/style");

const getLocalItem = require("../lib/storage").getLocalItem;

module.exports = function(state, emit) {
  state.username = getLocalItem("username");
  if (!state.username) {
    emit("pushState", "/");
    emit("render");
  }

  return html`
    <body class=${style}>
      <div class="welcome">
        <h1><span id="logo">_</span> ${state.username}</h1>
      </div>
    </body>
  `;
};
