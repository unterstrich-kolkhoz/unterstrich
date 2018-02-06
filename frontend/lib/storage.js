if (window && window.localStorage) {
  module.exports = {
    getLocalItem: localStorage.getItem.bind(localStorage),
    setLocalItem: localStorage.setItem.bind(localStorage)
  };
} else {
  module.exports = {
    getLocalItem: function() {},
    setLocalItem: function() {}
  };
}
