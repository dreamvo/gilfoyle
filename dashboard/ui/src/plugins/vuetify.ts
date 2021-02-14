import Vue from "vue";
import Vuetify from "vuetify/lib";
import colors from "vuetify/es5/util/colors";

Vue.use(Vuetify);

export default new Vuetify({
  theme: {
    dark: false,
    themes: {
      light: {
        primary: "#66f",
        secondary: "#34495e",
        accent: "#212121",
        info: "#34495e",
        warning: colors.orange.accent3,
        error: colors.red.accent4,
        success: colors.green.accent4
      }
    }
  }
});
