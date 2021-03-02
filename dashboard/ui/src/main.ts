import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import vuetify from "./plugins/vuetify";
import "./plugins/vee-validate";
import dayjs from "dayjs";

Vue.config.productionTip = false;

Vue.filter("readableDate", (date: string): string => {
  return dayjs(date).format("MMMM d, YYYY");
});

Vue.filter("readableDateHour", (date: string): string => {
  const fulldate = dayjs(date).format("MMMM d, YYYY");
  const hour = dayjs(date).format("h:mm a");

  return `${fulldate} at ${hour}`;
});

Vue.filter("readableNumber", (x: number | string): string => {
  if (!x) {
    return "0";
  }

  if (typeof x === "number") {
    x = x.toFixed(2).toString();
  }

  return x.replace(/\B(?<!\.\d*)(?=(\d{3})+(?!\d))/g, " ");
});

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount("#app");
