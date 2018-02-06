const choo = require("choo");
const app = choo();

app.use(require("./stores/login"));
app.route("/", require("./pages/home"));
app.route("/login", require("./pages/login"));
app.route("/signup", require("./pages/signup"));
app.route("/:username", require("./pages/user"));

if (module.parent) {
  module.exports = app;
} else {
  app.mount("body");
}
