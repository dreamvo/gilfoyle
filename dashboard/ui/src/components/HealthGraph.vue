<template>
  <v-card outlined :loading="loading">
    <v-card-title>{{ title }}</v-card-title>
    <v-card-subtitle v-if="avgResponseTime > 0 && $store.state.healthy"
      >Average response time: {{ avgResponseTime }} ms</v-card-subtitle
    >
    <v-system-bar height="54">
      <v-icon size="8" :color="$store.state.healthy ? 'green' : 'red'"
        >mdi-circle</v-icon
      >
      <span v-if="$store.state.healthy">Instance status: Running</span>
      <span v-else>Instance status: Unavailable</span>
    </v-system-bar>
  </v-card>
</template>

<script lang="ts">
import Vue from "vue";
import store from "../store";
import config from "../config";

interface Data {
  title: string;
  loading: boolean;
  avgResponseTime: number;
}

export default Vue.extend({
  name: "HealthGraph",
  data: (): Data => ({
    title: "Health check",
    loading: true,
    avgResponseTime: 0
  }),
  methods: {
    getAvgResponseTime(): number {
      if (!store.state.healthResponseTimes.length) {
        return 0;
      }
      const input: number[] = store.state.healthResponseTimes;
      this.avgResponseTime = Math.round(
        input.reduce((i, j) => i + j) / input.length
      );
      return this.avgResponseTime;
    }
  },
  created() {
    this.getAvgResponseTime();

    setInterval(this.getAvgResponseTime, config.healthCheckDelaySeconds * 1000);

    this.loading = false;
  }
});
</script>
