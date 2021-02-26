<template>
  <v-row>
    <v-col cols="12" md="12">
      <v-card outlined :loading="loading">
        <v-card-title>About this instance</v-card-title>

        <v-divider></v-divider>

        <v-card-text>
          <div>
            <h3 class="mb-2">Build</h3>

            <p>
              Gilfoyle version <strong>{{ version }}</strong
              >, commit <strong>{{ commit }}</strong>
            </p>

            <h3 class="mb-2">Configuration</h3>

            <div>Environment : {{ debug ? "Debug" : "Production" }}</div>
            <div>Database dialect : {{ dbDialect }}</div>
            <div>Storage driver : {{ storageDriver }}</div>
            <div>
              Max. file size : {{ (maxFileSize / 1024 / 1024).toFixed(2) }} MiB
            </div>
          </div>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";
import axios from "../services/axios";

interface HealthCheckResponse {
  commit: string;
  tag: string;
  database_dialect: string;
  debug: boolean;
  max_file_size: number;
  storage_driver: string;
}

interface Data {
  loading: boolean;
  commit: string;
  version: string;
  dbDialect: string;
  debug: boolean;
  maxFileSize: number;
  storageDriver: string;
}

export default Vue.extend({
  data: (): Data => ({
    loading: true,
    commit: "unknown",
    version: "unknown",
    dbDialect: "",
    debug: false,
    maxFileSize: 0,
    storageDriver: "unknown"
  }),
  async created() {
    const res = await axios.get<HealthCheckResponse>("/healthz");

    this.commit = res.data.commit;
    this.version = res.data.tag;
    this.dbDialect = res.data.database_dialect;
    this.debug = res.data.debug;
    this.maxFileSize = res.data.max_file_size;
    this.storageDriver = res.data.storage_driver;

    this.loading = false;
  }
});
</script>
