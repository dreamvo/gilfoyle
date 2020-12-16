import Vue from "vue";
import Vuex from "vuex";
import { AxiosResponse } from "axios";
import axios from "../services/axios";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    healthy: false,
    healthResponseTimes: [] as number[]
  },
  mutations: {
    setHealthStatus(state, status: boolean) {
      state.healthy = status;
    },
    addResponseTime(state, d: number) {
      state.healthResponseTimes.push(d);
    },
    resetResponseTime(state) {
      state.healthResponseTimes = [];
    }
  },
  actions: {
    async healthCheck(context) {
      const t1: number = Date.now();
      const res: AxiosResponse | void = await axios
        .get("/healthz")
        .catch((/*err: Error*/) => {
          context.commit("setHealthStatus", false);
          context.commit("resetResponseTime");
        });

      if (res && res.status == 200) {
        context.commit("setHealthStatus", true);
      }

      if (context.state.healthy) {
        context.commit("addResponseTime", Date.now() - t1);
      }
    }
  },
  modules: {}
});
