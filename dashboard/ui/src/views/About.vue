<template>
  <v-row>
    <v-col cols="12" md="12">
      <v-card outlined :loading="loading">
        <v-card-title>About this instance</v-card-title>

        <v-divider></v-divider>

        <v-card-text>
          <div>
            <h3>Build</h3>

            <p class="mt-2">Gilfoyle version <strong>{{version}}</strong>, commit <strong>{{commit}}</strong></p>
          </div>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";
import MediaForm from "../components/MediaForm.vue";
import axios from "../services/axios";

interface Data {
  loading: boolean
  commit: string
  version: string
}

export default Vue.extend({
  components: {MediaForm},
  data: (): Data => ({
    loading: true,
    commit: 'unknown',
    version: 'unknown'
  }),
  methods: {},
  async created() {
    const res = await axios.get('/healthz')

    this.commit = res.data.commit
    this.version = res.data.tag

    this.loading = false
  }
});
</script>
