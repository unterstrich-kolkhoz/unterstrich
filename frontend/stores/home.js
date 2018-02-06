module.exports = function(state, emitter) {
  state.title = "";

  emitter.on("update", function({ value }) {
    state.title = value;
    emitter.emit("render");
  });
};
